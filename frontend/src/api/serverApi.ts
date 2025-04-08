import { MessageData } from "@/components/Message";

const BASE_URL = "/api";

interface RawMessageData {
    messageid: number;
    channelid: number;
    userid: number;
    username?: string;
    date: Date;
    message: string;
}

interface RawChannelData {
    ChannelId: number;
    ChannelName: string;
}

interface ServerIconResponse {
    ServerId: number;
    ServerName: string;
    ImageUrl: string | undefined;
}

export const fetchServerMessages = async (
    serverId: number,
    messageCount: number
): Promise<MessageData[]> => {
    const response = await fetch(`${BASE_URL}/servers/${serverId}/messages?count=${messageCount}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
    });

    if (!response.ok) {
        throw new Error(`Failed to fetch messages: ${response.statusText}`);
    }

    const data = await response.json();
    const messageDataArray: MessageData[] = data["messages"].map((msg: RawMessageData) => {
        const message: MessageData = {
            message_id: msg.messageid,
            channel_id: msg.channelid,

            author: msg.username ? msg.username : "User" + msg.userid,
            author_id: msg.userid.toString(),
            date: new Date(msg.date),
            message: msg.message,
        };

        return message;
    });

    return messageDataArray.sort((a, b) => a.message_id - b.message_id);
};

export const fetchChannels = async (serverId: number): Promise<RawChannelData[]> => {
    const response = await fetch(`${BASE_URL}/servers/${serverId}/channels`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
    });

    if (!response.ok) {
        throw new Error(`Failed to fetch channels: ${response.statusText}`);
    }

    const data = await response.json();
    const channelInfoArray: RawChannelData[] = data["channels"];
    return channelInfoArray;
};

export const fetchUserServers = async (userId: number): Promise<ServerIconResponse[]> => {
    const response = await fetch(`${BASE_URL}/users/${userId}/servers`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
    });

    if (!response.ok) {
        throw new Error(`Failed to fetch servers: ${response.statusText}`);
    }

    const data = await response.json();
    return data.servers.map((server: ServerIconResponse) => ({
        ServerId: server.ServerId,
        ServerName: server.ServerName,
        ImageUrl: "https://miro.medium.com/v2/resize:fit:720/format:webp/0*UD_CsUBIvEDoVwzc.png",
    }));
};
