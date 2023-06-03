const net = require('net');

function numToUint8Array(num) {
    let arr = new Uint8Array(8);
  
    for (let i = 0; i < 8; i++) {
      arr[i] = num % 256;
      num = Math.floor(num / 256);
    }
  
    return arr;
}

class Client {
    constructor() {
        this._socket = new net.Socket();
        this._socket.connect(6789, '127.0.0.1', () => {
            this._socket.write(Uint8Array.from([1, 10, ...Uint8Array.from(Buffer.from('testing123'))]));
        });

        this._socket.on('data', buf => {
            console.info(buf);
        });

        this._socket.on('end', () => {
            this._socket.destroy();
        });
    }
}

new Client();
