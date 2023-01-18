const net = require('net');

const encoder = new TextEncoder();
const socket = new net.Socket();

socket.connect(9876, 'localhost', () => {
    console.info('connected...');
    socket.write('<1::>');
});

socket.on('data', (frame) => {
    console.info(frame);
    console.info(frame.toString());
});

socket.on('close', () => {
    socket.destroy();
});
