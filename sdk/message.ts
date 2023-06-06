interface MessageTypePayload {
    readonly connect: {
        readonly clientId: string;
        readonly username: string;
        readonly password: string;
    };
    readonly connectAck: {
        readonly sessionId: string;
    };
}

const MESSAGE_TYPE_CODE = {
    connect: 1,
    connectAck: 2
};

const CODE_MESSAGE_TYPE = {
    1: 'connect',
    2: 'connectAck'
};

const MESSAGE_TYPE_TRANSFORM = {
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
    }
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