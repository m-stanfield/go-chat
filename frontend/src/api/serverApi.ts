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
    ImageUrl: string | undefined;
}

export const fetchServerMessages = async (
    serverId: number,
    messageCount: number
): Promise<RawMessageData[]> => {
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
    const messageDataArray: RawMessageData[] = data["messages"].map((msg: RawMessageData) => {
        msg.username = msg.username ? msg.username : "User" + msg.userid;
        return msg;
    });

    return messageDataArray.sort((a, b) => a.messageid - b.messageid);
};

export const fetchChannels = async (serverId: number): Promise<RawChannelData[]> => {
    const response = await fetch(`http://localhost:8080/api/servers/${serverId}/channels`, {
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
    const response = await fetch(`http://localhost:8080/api/users/${userId}/servers`, {
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
