import { MessageData } from "@/components/Message";
import { Channel } from "@/types/channel";

const BASE_URL = "/api";

interface RawMessageData {
    messageid: number;
    channelid: number;
    serverid: number;
    userid: number;
    date: Date;
    message: string;
}

interface RawChannelData {
    ChannelId: number;
    ServerId: number;
    ChannelName: string;
    Timestamp: Date;
}

export interface ServerIconResponse {
    ServerId: number;
    ServerName: string;
    ImageUrl: string | undefined;
}

export interface CreateChannelResponse {
    channelId: number;
}

export const postChannel = async (
    serverId: number,
    channelName: string
): Promise<CreateChannelResponse> => {
    const response = await fetch(`${BASE_URL}/servers/${serverId}/channels`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ channelName }),
    });

    if (!response.ok) {
        throw new Error(`Failed to create channel: ${response.statusText}`);
    }
    const data = await response.json();
    if (!data["channelid"]) {
        throw new Error("Invalid response from server");
    }
    const channelId: CreateChannelResponse = {
        channelId: data["channelid"],
    };
    return channelId;
}

export const postServers = async (
    serverName: string,
    imageUrl: string
): Promise<Response> => {
    const response = await fetch(`${BASE_URL}/servers`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ serverName, imageUrl }),
    });

    if (!response.ok) {
        throw new Error(`Failed to create server: ${response.statusText}`);
    }
    return response;
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
    if (!data["messages"]) {
        return [];
    }
    const messageDataArray: MessageData[] = data["messages"].map((msg: RawMessageData) => {
        const message: MessageData = {
            message_id: msg.messageid,
            channel_id: msg.channelid,
            server_id: msg.serverid,
            author: "User" + msg.userid,
            author_id: msg.userid,
            date: new Date(msg.date),
            message: msg.message,
        };

        return message;
    });

    return messageDataArray.sort((a, b) => a.message_id - b.message_id);
};

export const fetchChannels = async (serverId: number): Promise<Channel[]> => {
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
    if (!data["channels"]) {
        return [];
    }
    const channelInfoArray: Channel[] = data["channels"].map((channel: RawChannelData) => {
        return {
            ChannelId: channel.ChannelId,
            ServerId: serverId,
            ChannelName: channel.ChannelName,
            Timestamp: channel.Timestamp,
        };
    });
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
