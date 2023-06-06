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

    constructor(private readonly _options: ClientOptions) {
        this._socket = new net.Socket();
        this._socket.on('data', this._onData.bind(this));
        this._socket.on('end', this.close.bind(this));
    }

    open(options: ClientConnectOptions) {
        return new Promise<void>((resolve, reject) => {
            this._socket.once('error', reject);
            this._socket.connect(options.port || 6789, options.host, () => {
                this._socket.write(new Message('connect', {
                    clientId: this._options.id,
                    username: options.username || 'admin',
                    password: options.password || 'admin'
                }).serialize(), err => {
                    if (err) {
                        return reject(err);
                    }

                    resolve();
                });
            });
        });
    }

    close() {
        if (this._socket.closed) {
            return;
        }

        this._socket.destroy();
    }

    send(payload: string) {
        return new Promise<void>((resolve, reject) => {
            this._socket.write(Buffer.from([
                3,
                ...Buffer.from([0, 0, 0, payload.length]),
                ...Buffer.from(payload)
            ]), err => {
                if (err) {
                    return reject(err);
                }

                resolve();
            });
        });
    }

    private _onData(buf: Buffer) {
        console.log(new Message(buf));
    }
}