import ChatPage from "./ChatPage";
import { SyntheticEvent, useEffect, useRef, useState } from "react";
import { MessageData } from "./Message";
import ChannelIconBanner, { ChannelInfo } from "./ChannelIconBanner";

interface ServerPageProps {
    server_id: number;
}
function ServerPage({ server_id }: ServerPageProps) {
    const [selectedChannelId, setSelectedChannelId] = useState<ChannelInfo>({
        channel_id: 0,
        channel_name: "",
    });
    const [channnelMessages, setChannelMessages] = useState<
        Map<number, MessageData[]>
    >(new Map());
    useEffect(() => {
        (async () => {
            try {
                // Send POST request to backend
                const response = await fetch(
                    `http://localhost:8080/api/server/${server_id}/messages`,
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
                    const newmap = new Map([
                        [messageDataArray[0].channel_id, messageDataArray],
                    ]);
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
    }, [server_id, selectedChannelId]);
    const ws = useRef<WebSocket | null>(null);
    const maxMessageLength = 30;
    const onSubmit = (t: SyntheticEvent, inputValue: string): string => {
        t.preventDefault();
        if (inputValue.length === 0) {
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
                    if (channel_messages.length > maxMessageLength) {
                        messages.set(channel_id, channel_messages.slice(-maxMessageLength));
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
    }, []);

    return (
        <div className="flex flex-grow bg-red-500 flex-col">
            <ChannelIconBanner
                server_id={server_id}
                onChannelSelect={setSelectedChannelId}
            />
            <ChatPage
                channel_id={selectedChannelId?.channel_id}
                onSubmit={onSubmit}
                messages={channnelMessages.get(selectedChannelId?.channel_id) || []}
            />
        </div>
    );
}

export default ServerPage;
