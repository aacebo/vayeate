const net = require('net');

const socket = net.connect(9876, '127.0.0.1', () => {
    socket.write('<1::>');
    socket.write('<3:test:>');
});

socket.on('data', buf => {
    console.log(buf.toString());
});

socket.on('end', () => {
    socket.destroy();
});
