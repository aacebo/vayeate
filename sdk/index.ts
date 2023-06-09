import { Client } from './client';

(async () => {
    const a = new Client({ id: 'a' });
    const b = new Client({ id: 'b' });

    await a.open({
        host: '127.0.0.1',
        port: 6789
    });

    await b.open({
        host: '127.0.0.1',
        port: 6789
    });

    let counter = 0;

    await a.subscribe('aacebo.test', _ => {
        counter++
        console.info('a', counter);
    });

    await b.subscribe('aacebo.test', _ => {
        counter++
        console.info('b', counter);
    });

    for (let i = 0; i < 500000; i++) {
        a.publish('aacebo.test', Buffer.from('testing123!'));
        b.publish('aacebo.test', Buffer.from('hello world!'));
    }
})();