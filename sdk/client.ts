import net from 'net';

import { connectMessage } from './message';

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
                this._socket.write(connectMessage(
                    this._options.id,
                    options.username || 'admin',
                    options.password || 'admin'
                ), err => {
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
        const code = buf.at(0);
        const length = buf.readUInt32BE(1);
        const payload = buf.subarray(5);
        const sessionIdLength = buf.readUInt32BE(5);
        const sessionId = buf.subarray(9);

        console.log('code', code);
        console.log('length', length);
        console.log('payload', payload);
        console.log('session_id length', sessionIdLength);
        console.log('session_id', sessionId.toString());
    }
}