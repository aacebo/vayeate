import { Client } from './client';

(async () => {
    const client = new Client({ id: 'a' });

    console.info(await client.open({
        host: '127.0.0.1',
        port: 6789
    }));

    await client.publish('aacebo.test', Buffer.from('testing123!'));
    await client.subscribe('aacebo.test', m => {
        console.info(m);
    });
})();