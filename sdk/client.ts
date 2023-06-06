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
    private readonly _socket: net.Socket;

    get sessionId() { return this._sessionId; }
    private _sessionId?: string;

    constructor(private readonly _options: ClientOptions) {
        this._socket = new net.Socket();
        this._socket.on('data', this._onData.bind(this));
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
                    const m = new Message<'connectAck'>(buf);

                    if (m.type !== 'connectAck') {
                        return reject(new Error('connect handshake incomplete'));
                    }

                    this._sessionId = m.payload.sessionId;
                    resolve(m.payload.sessionId);
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

        if (this._socket.closed) {
            return;
        }

        this._socket.destroy();
    }

    publish(topic: string, payload: Buffer) {
        return new Promise<void>((resolve, reject) => {
            const m = new Message('publish', { topic, payload });

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

    private _onData(buf: Buffer) {
        console.log(new Message(buf));
    }
}