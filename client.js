const net = require('net');

for (let i = 0; i < 100; i++) {
    const socket = net.connect(9876, '127.0.0.1', () => {
        socket.write('<1::>'); // ping
        socket.write('<3:test:>'); // assert
        socket.write('<5:test:>'); // consume
        socket.write('<4:test:my test message>'); // produce
    });

    socket.on('data', buf => {
        console.log(buf.toString());
    });

    socket.on('end', () => {
        socket.destroy();
    });
}
