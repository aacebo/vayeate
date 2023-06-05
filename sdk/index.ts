import { Client } from './client';

(async () => {
    const client = new Client();
    await client.open('127.0.0.1', 6789);
    await client.send('testing123!');
})();