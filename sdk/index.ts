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
    })

    await a.subscribe('aacebo.test', _ => {
        console.info('a');
    });

    await b.subscribe('aacebo.test', _ => {
        console.info('b');
    });

    for (let i = 0; i < 20; i++) {
        await a.publish('aacebo.test', Buffer.from('testing123!'));
    }
})();