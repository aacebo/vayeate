interface MessageTypePayload {
    readonly error: {
        readonly reason: string;
    };
    readonly connect: {
        readonly clientId: string;
        readonly username: string;
        readonly password: string;
    };
    readonly connectAck: {
        readonly sessionId: string;
    };
    readonly publish: {
        readonly topic: string;
        readonly payload: Buffer;
    };
    readonly publishAck: { };
    readonly ping: { };
    readonly pingAck: { };
}

const MESSAGE_TYPE_CODE = {
    error: 0,
    connect: 1,
    connectAck: 2,
    publish: 3,
    publishAck: 4,
    ping: 12,
    pingAck: 13
};

const CODE_MESSAGE_TYPE = {
    0: 'error',
    1: 'connect',
    2: 'connectAck',
    3: 'publish',
    4: 'publishAck',
    12: 'ping',
    13: 'pingAck'
};

const MESSAGE_TYPE_TRANSFORM = {
    error: (b: Buffer) => {
        const len = b.readUint32BE(5);
        const value = b.subarray(9, 9 + len);

        return {
            reason: value.toString()
        };
    },
    connect: (b: Buffer) => {
        let i = 5;
        let len = b.readUint32BE(i);
        const clientId = b.subarray(i + 4, i + 4 + len);
        i = i + 4 + len;

        len = b.readUint32BE(i);
        const username = b.subarray(i + 4, i + 4 + len);
        i = i + 4 + len;
        
        len = b.readUInt32BE(i);
        const password = b.subarray(i + 4, i + 4 + len);

        return {
            clientId: clientId.toString(),
            username: username.toString(),
            password: password.toString()
        };
    },
    connectAck: (b: Buffer) => {
        const len = b.readUint32BE(5);
        const value = b.subarray(9, 9 + len);
        
        return {
            sessionId: value.toString()
        };
    },
    publish: (b: Buffer) => {
        let i = 5;
        let len = b.readUint32BE(i);
        const topic = b.subarray(i + 4, i + 4 + len);
        i = i + 4 + len;

        len = b.readUint32BE(i);
        const payload = b.subarray(i + 4, i + 4 + len);

        return {
            topic: topic.toString(),
            payload
        };
    },
    publishAck: (_: Buffer) => ({ }),
    ping: (_: Buffer) => ({ }),
    pingAck: (_: Buffer) => ({ })
};

export class Message<T extends keyof MessageTypePayload> {
    readonly type: T;
    readonly payload: MessageTypePayload[T];

    constructor(type: T, payload: MessageTypePayload[T])
    constructor(buf: Buffer)
    constructor(...args: any[]) {
        if (args.length == 2) {
            this.type = args[0];
            this.payload = args[1];
        } else {
            const buf: Buffer = args[0];
            this.type = CODE_MESSAGE_TYPE[buf.at(0)!];
            this.payload = MESSAGE_TYPE_TRANSFORM[this.type](buf) as any;
        }
    }

    serialize() {
        const code = Buffer.from([MESSAGE_TYPE_CODE[this.type]]);
        let payload = Buffer.from([ ]);
        
        for (const key in this.payload) {
            const length = Buffer.from([0, 0, 0, 0]);
            length.writeUInt32BE((this.payload[key] as string).length);

            payload = Buffer.concat([
                payload,
                length,
                Buffer.from(this.payload[key] as string)
            ]);
        }

        const length = Buffer.from([0, 0, 0, 0]);
        length.writeUInt32BE(payload.length);

        return Buffer.concat([
            code,
            length,
            payload
        ]);
    }
}