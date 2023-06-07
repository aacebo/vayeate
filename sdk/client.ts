import net from 'net';

import { Message } from './message';

export interface ClientOptions {
    readonly id: string;
}

export interface ClientConnectOptions {
    readonly host: string;
    readonly port?: number;
    readonly username?: string;
    readonly password?: string;
}

export class Client {
    get sessionId() { return this._sessionId; }
    private _sessionId?: string;

    private readonly _socket: net.Socket;
    private _pingTimer?: NodeJS.Timer;
    private readonly _subscriptions: { [topic: string]: (message: Message<'consume'>) => void } = { };

    constructor(private readonly _options: ClientOptions) {
        this._socket = new net.Socket();
        this._socket.on('end', this.close.bind(this));
    }

    /**
     * open a new connection and authenticate
     * @param options
     * @returns Session ID
     */
    open(options: ClientConnectOptions) {
        return new Promise<string>((resolve, reject) => {
            this._socket.once('error', reject);
            this._socket.connect(options.port || 6789, options.host, () => {
                this._socket.once('data', buf => {
                    const m = new Message(buf);

                    if (m.type !== 'connectAck') {
                        let message = 'connect handshake incomplete';

                        if (m.type === 'error') {
                            message = (m as Message<'error'>).payload.reason;
                        }

                        return reject(new Error(message));
                    }

                    this._onConncted();
                    this._sessionId = (m as Message<'connectAck'>).payload.sessionId;
                    resolve(this._sessionId);
                });

                this._socket.write(new Message('connect', {
                    clientId: this._options.id,
                    username: options.username || 'admin',
                    password: options.password || 'admin'
                }).serialize(), err => {
                    if (err) {
                        return reject(err);
                    }
                });
            });
        });
    }

    close() {
        this._sessionId = undefined;
        clearInterval(this._pingTimer);
        this._pingTimer = undefined;

        if (this._socket.closed) {
            return;
        }

        this._socket.destroy();
    }

    publish(topic: string, payload: Buffer) {
        return new Promise<void>((resolve, reject) => {
            const m = new Message('publish', { topic, payload });

            this._socket.once('error', reject);
            this._socket.once('data', buf => {
                const ack = new Message<'publishAck'>(buf);

                if (ack.type !== 'publishAck') {
                    return reject(new Error('waiting for publish acknowledgement'));
                }

                resolve();
            });

            this._socket.write(m.serialize(), err => {
                if (err) {
                    return reject(err);
                }
            });
        });
    }

    subscribe(topic: string, cb: (message: Message<'consume'>) => void) {
        return new Promise<void>((resolve, reject) => {
            const m = new Message('subscribe', { topic });

            this._socket.once('error', reject);
            this._socket.once('data', buf => {
                const ack = new Message<'subscribeAck'>(buf);

                if (ack.type !== 'subscribeAck') {
                    return reject(new Error('waiting for subscribe acknowledgement'));
                }

                this._subscriptions[topic] = cb;
                resolve();
            });

            this._socket.write(m.serialize(), err => {
                if (err) {
                    return reject(err);
                }
            });
        });
    }

    private _onConncted() {
        this._socket.on('data', this._onData.bind(this));
        this._pingTimer = setInterval(() => {
            this._socket.write(new Message('ping', { }).serialize());
        }, 30000);
    }

    private _onData(buf: Buffer) {
        const m = new Message(buf);

        if (m.type == 'consume' && this._subscriptions[(m as Message<'consume'>).payload.topic]) {
            this._subscriptions[(m as Message<'consume'>).payload.topic](m as Message<'consume'>);
            this._socket.write(new Message('consumeAck', {
                topic: (m as Message<'consume'>).payload.topic
            }).serialize());
        }
    }
}