import net from 'net';

export class Client {
    private readonly _socket: net.Socket;

    constructor() {
        this._socket = new net.Socket();
        this._socket.on('data', this._onData.bind(this));
        this._socket.on('end', this.close.bind(this));
    }

    open(host: string, port: number) {
        return new Promise<void>((resolve, reject) => {
            this._socket.once('error', reject);
            this._socket.connect(port, host, resolve);
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
        console.info(buf);
    }
}