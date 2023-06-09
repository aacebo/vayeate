import net from 'net';

import { CODE_MESSAGE_TYPE, MESSAGE_TYPE_CODE, MESSAGE_TYPE_TRANSFORM, Message } from './message';

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
    private readonly _socket: net.Socket;
    private _pingTimer?: NodeJS.Timer;
    private readonly _subscriptions: { [topic: string]: (message: Message<'consume'>) => void } = { };

    private _buffer = Buffer.from([ ]);

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
        return new Promise<void>((resolve, reject) => {
            this._socket.once('error', reject);
            this._socket.connect(options.port || 6789, options.host, () => {
                this._socket.once('data', buf => {
                    const code = buf.readUInt8();

                    if (!code || code != MESSAGE_TYPE_CODE.connectAck) {
                        let message = 'connect handshake incomplete';
                        return reject(new Error(message));
                    }

                    this._onConncted();
                    resolve();
                });

                this._socket.write(new Message('connect', Date.now(), {
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
        clearInterval(this._pingTimer);
        this._pingTimer = undefined;

        if (this._socket.closed) {
            return;
        }

        this._socket.removeAllListeners();
        this._socket.destroy();
    }

    publish(topic: string, payload: Buffer) {
        return new Promise<void>((resolve, reject) => {
            const m = new Message('publish', Date.now(), { topic, payload });

            this._socket.write(m.serialize(), err => {
                if (err) {
                    return reject(err);
                }

                resolve();
            });
        });
    }

    subscribe(topic: string, cb: (message: Message<'consume'>) => void) {
        return new Promise<void>((resolve, reject) => {
            const m = new Message('subscribe', Date.now(), { topic });

            this._socket.write(m.serialize(), err => {
                if (err) {
                    return reject(err);
                }

                this._subscriptions[topic] = cb;
                resolve();
            });
        });
    }

    private _onConncted() {
        this._socket.on('data', this._onData.bind(this));
        this._pingTimer = setInterval(() => {
            this._socket.write(new Message('ping', Date.now(), { }).serialize());
        }, 30000);
    }

    private _onData(b: Buffer) {
        this._buffer = Buffer.concat([this._buffer, b]);

        while (this._buffer.length >= 13) {
            const code = this._buffer.readUint8();
            const sentAt = this._buffer.readBigUInt64BE(1);
            const length = this._buffer.readUint32BE(9);

            // incomplete packet
            if (13 + length > this._buffer.length) {
                break;
            }

            const type = CODE_MESSAGE_TYPE[code];
            const payload = MESSAGE_TYPE_TRANSFORM[type](this._buffer.subarray(13, 13 + length));
            const m = new Message(type, sentAt, payload);

            if (m.type === 'consume') {
                if (this._subscriptions[m.payload.topic]) {
                    this._subscriptions[m.payload.topic](m);
                }

                this._socket.write(
                    new Message('consumeAck', Date.now(), {
                        topic: m.payload.topic
                    }).serialize()
                );
            }

            this._buffer = this._buffer.subarray(13 + length);
        }
    }
}
