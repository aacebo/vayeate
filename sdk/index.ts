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

    await b.subscribe('aacebo.test', m => {
        console.info(m);
    });

    await a.publish('aacebo.test', Buffer.from('testing123!'));
})();