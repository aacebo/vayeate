export enum MessageType {
    Connect = 1,
    ConnectAck = 2,
    Publish = 3,
    PublishAck = 4,
    PublishRec = 5,
    PublishRel = 6,
    PublishComp = 7,
    Subscribe = 8,
    SubscribeAck = 9,
    Unsubscribe = 10,
    UnsubscribeAck = 11,
    Ping = 12,
    PingAck = 13,
    Disconnect = 14
}

export function connectMessage(clientId: string, username: string, password: string) {
    const payload = Buffer.from([
        ...Buffer.from([0, clientId.length]),
        ...Buffer.from(clientId),
        ...Buffer.from([0, username.length]),
        ...Buffer.from(username),
        ...Buffer.from([0, password.length]),
        ...Buffer.from(password)
    ]);

    return Buffer.from([
        MessageType.Connect,
        ...Buffer.from([0, 0, 0, payload.length]),
        ...payload
    ]);
}