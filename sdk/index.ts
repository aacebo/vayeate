import { Client } from './client';

(async () => {
    const client = new Client({ id: 'a' });

    await client.open({
        host: '127.0.0.1',
        port: 6789
    });

    await client.send('testing123!');
})();