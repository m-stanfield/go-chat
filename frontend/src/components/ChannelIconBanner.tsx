import React, { useEffect, useState } from "react";

export type ChannelInfo = {
    channel_id: number;
    channel_name: string;
};
interface ChannelIconBannerProps {
    server_id: number;
    onChannelSelect: (channel_id: ChannelInfo) => void;
}
function ChannelIconBanner({
    server_id,
    onChannelSelect,
}: ChannelIconBannerProps) {
    const [channelInformationArray, setChannelInformationArray] = useState<
        ChannelInfo[]
    >([]);

    function onChannelSelectGenerator(channel_id: ChannelInfo) {
        return (t: React.MouseEvent<HTMLElement>) => {
            t.preventDefault();
            onChannelSelect(channel_id);
        };
    }
    useEffect(() => {
        (async () => {
            try {
                // Send POST request to backend
                const response = await fetch(
                    `http://localhost:8080/api/server/${server_id}/channels`,
                    {
                        method: "GET",
                        headers: {
                            "Content-Type": "application/json",
                        },
                        credentials: "include",
                    },
                );

                // Handle response
                if (response.ok) {
                    const data = await response.json();
                    const channelInfoArray: ChannelInfo[] = data["channels"].map(
                        (msg) => {
                            const obj = {
                                channel_id: msg.ChannelId,
                                channel_name: msg.ChannelName,
                            };
                            return obj;
                        },
                    );
                    channelInfoArray.sort((a, b) => a.channel_id - b.channel_id);
                    setChannelInformationArray(() => {
                        return channelInfoArray;
                    });
                    if (channelInfoArray.length > 0) {
                        onChannelSelect(channelInfoArray[0]);
                    }
                } else {
                    console.error("Login failed:", response.statusText);
                    setChannelInformationArray(() => {
                        return [];
                    });

                    return;
                }
            } catch (error) {
                console.error("Error submitting login:", error);
                return;
            }
        })();
    }, [server_id, onChannelSelect]);

    const items = channelInformationArray.map((s) => (
        <div key={s.channel_id} className="">
            <button onClick={onChannelSelectGenerator(s)}>
                Channel: {s.channel_id}
            </button>
        </div>
    ));
    return <div>{items}</div>;
}

export default ChannelIconBanner;
