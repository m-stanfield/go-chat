import ChatPage from "./ChatPage";
import { SyntheticEvent, useEffect, useRef, useState } from "react";
import { MessageData } from "./Message";
import IconBanner, { IconInfo } from "./IconList";

interface ServerPageProps {
    server_id: number;
    number_of_messages: number;
}
function ServerPage({ server_id, number_of_messages }: ServerPageProps) {
    const [selectedChannelId, setSelectedChannelId] = useState<number>(0);
    const [channnelMessages, setChannelMessages] = useState<
        Map<number, MessageData[]>
    >(new Map());
    useEffect(() => {
        (async () => {
            if (server_id < 0) {
                return;
            }
            try {
                // Send POST request to backend
                const response = await fetch(
                    `http://localhost:8080/api/server/${server_id}/messages?count=${number_of_messages}`,
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
                    const messageDataArray: MessageData[] = data["messages"].map(
                        (msg) => {
                            const obj = {
                                message_id: msg.messageid,
                                channel_id: msg.channelid,
                                author: msg.username,
                                author_id: msg.userid,
                                date: new Date(msg.date), // Convert to JavaScript Date object
                                message: msg.message,
                            };
                            return obj;
                        },
                    );
                    messageDataArray.sort((a, b) => a.message_id - b.message_id);
                    console.log(data);
                    const newmap = messageDataArray.reduce((map, obj) => {
                        const { channel_id } = obj;
                        if (!map.has(channel_id)) {
                            map.set(channel_id, []);
                        }
                        map.get(channel_id).push(obj);
                        return map;
                    }, new Map());
                    setChannelMessages(() => {
                        return newmap;
                    });
                } else {
                    console.error("Login failed:", response.statusText);
                    setChannelMessages(() => {
                        return new Map();
                    });

                    return;
                }
            } catch (error) {
                console.error("Error submitting login:", error);
                return;
            }
        })();
    }, [server_id, selectedChannelId, number_of_messages]);
    const ws = useRef<WebSocket | null>(null);
    const onSubmit = (t: SyntheticEvent, inputValue: string): string => {
        t.preventDefault();
        if (inputValue.length === 0) {
            return inputValue;
        } else if (inputValue.length >= 1000) {
            return inputValue;
        }
        const stringified = JSON.stringify({
            channel_id: selectedChannelId,
            message: inputValue,
        });
        if (ws === null) {
            console.log("websocket hasn't be initialized yet");
            return inputValue;
        }
        if (ws.current?.readyState === WebSocket.CLOSED) {
            console.log("can't send ws closed");
            return inputValue;
        }
        ws.current?.send(stringified);
        return "";
    };

    useEffect(() => {
        const startTime = Date.now();
        console.log("starting to open websocket", Date.now() - startTime);
        ws.current = new WebSocket(`ws://localhost:8080/websocket`);
        ws.current.onopen = () => {
            console.log("opening ws", Date.now() - startTime);
        };
        ws.current.onmessage = function(event: MessageEvent) {
            const json = JSON.parse(event.data);
            try {
                const newMessage: MessageData = {
                    message_id: json.messageid,
                    channel_id: json.channelid,
                    message: json.message,
                    date: new Date(json.date),
                    author: json.username,
                    author_id: json.userid,
                };
                const channel_id = newMessage.channel_id;
                if (!channel_id) {
                    return;
                }
                setChannelMessages((messages) => {
                    let channel_messages = messages.get(channel_id);
                    if (!channel_messages) {
                        return messages;
                    }
                    channel_messages = [...channel_messages, newMessage];
                    if (channel_messages.length > number_of_messages) {
                        messages.set(
                            channel_id,
                            channel_messages.slice(-number_of_messages),
                        );
                    } else {
                        messages.set(channel_id, channel_messages);
                    }
                    return new Map(messages);
                });
            } catch (err) {
                console.log(err);
            }
        };
        ws.current.onclose = () => {
            console.log("ws closed", Date.now() - startTime);
        };

        return () => {
            ws.current?.close();
        };
    }, [number_of_messages]);

    const [channelInformationArray, setChannelInformationArray] = useState<
        IconInfo[]
    >([]);

    useEffect(() => {
        (async () => {
            if (server_id < 0) {
                return;
            }
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
                    const channelInfoArray: IconInfo[] = data["channels"].map((msg) => {
                        const obj: IconInfo = {
                            icon_id: msg.ChannelId,
                            name: msg.ChannelName,
                            image_url: undefined,
                        };
                        return obj;
                    });
                    channelInfoArray.sort((a, b) => a.icon_id - b.icon_id);
                    setChannelInformationArray(() => {
                        return channelInfoArray;
                    });
                    if (channelInfoArray.length > 0) {
                        setSelectedChannelId(channelInfoArray[0].icon_id);
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
    }, [server_id]);
    return (
        <div className="flex flex-col flex-grow overflow-y-auto">
            <IconBanner
                icon_info={channelInformationArray}
                onServerSelect={setSelectedChannelId}
            />
            <ChatPage
                channel_id={selectedChannelId}
                onSubmit={onSubmit}
                messages={channnelMessages.get(selectedChannelId) || []}
            />
        </div>
    );
}

export default ServerPage;
