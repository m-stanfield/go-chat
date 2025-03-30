import { MessageData } from "../components/Message";
import { IconInfo } from "../components/IconList";

interface RawMessageData {
    messageid: number;
    channelid: number;
    userid: number;
    username?: string;
    date: string;
    message: string;
}

interface RawChannelData {
    ChannelId: number;
    ChannelName: string;
}

interface ServerIconResponse {
    ServerId: number;
    ServerName: string;
    image_url: string | undefined;
}

export const fetchServerMessages = async (serverId: number, messageCount: number): Promise<MessageData[]> => {
    const response = await fetch(
        `http://localhost:8080/api/servers/${serverId}/messages?count=${messageCount}`,
        {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
            credentials: "include",
        }
    );

    if (!response.ok) {
        throw new Error(`Failed to fetch messages: ${response.statusText}`);
    }

    const data = await response.json();
    const messageDataArray: MessageData[] = data["messages"].map((msg: RawMessageData) => {
        const username = msg.username ? msg.username : "User" + msg.userid;
        return {
            message_id: msg.messageid,
            channel_id: msg.channelid,
            author: username,
            author_id: msg.userid,
            date: new Date(msg.date),
            message: msg.message,
        };
    });

    return messageDataArray.sort((a, b) => a.message_id - b.message_id);
};

export const fetchChannels = async (serverId: number): Promise<IconInfo[]> => {
    const response = await fetch(
        `http://localhost:8080/api/servers/${serverId}/channels`,
        {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
            credentials: "include",
        }
    );

    if (!response.ok) {
        throw new Error(`Failed to fetch channels: ${response.statusText}`);
    }

    const data = await response.json();
    const channelInfoArray: IconInfo[] = data["channels"].map((msg: RawChannelData) => ({
        icon_id: msg.ChannelId,
        name: msg.ChannelName,
        image_url: undefined,
    }));

    return channelInfoArray.sort((a, b) => a.icon_id - b.icon_id);
};

export const fetchUserServers = async (userId: number): Promise<IconInfo[]> => {
    const response = await fetch(
        `http://localhost:8080/api/users/${userId}/servers`,
        {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
            credentials: "include",
        }
    );

    if (!response.ok) {
        throw new Error(`Failed to fetch servers: ${response.statusText}`);
    }

    const data = await response.json();
    return data.servers.map((server: ServerIconResponse) => ({
        icon_id: server.ServerId,
        name: server.ServerName,
        image_url: "https://miro.medium.com/v2/resize:fit:720/format:webp/0*UD_CsUBIvEDoVwzc.png",
    }));
}; 