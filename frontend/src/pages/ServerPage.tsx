import { SyntheticEvent, useEffect, useRef, useState } from "react";
import { fetchServerMessages, fetchChannels } from "../api/serverApi";
import ChatPage from "@/components/ChatWindow";
import { MessageData } from "@/components/Message";
import ChannelSidebar from "@/components/ChannelSidebar";
import { Channel } from "@/types/channel";
import { useAuth } from "@/AuthContext";
import { toast } from "sonner";
import { useNavigate, useParams } from "react-router-dom";

interface ServerPageProps {
    server_id: number;
    number_of_messages: number;
}

function ServerPage({ server_id, number_of_messages }: ServerPageProps) {
    const auth = useAuth();
    const navigate = useNavigate();
    const { channelId: channelIdStr } = useParams<{ channelId: string }>();
    const channelId = channelIdStr ? parseInt(channelIdStr) : -1;

    const [channels, setChannels] = useState<Channel[]>([]);
    const [channnelMessages, setChannelMessages] = useState<Map<number, MessageData[]>>(new Map());
    useEffect(() => {
        (async () => {
            if (server_id < 0 || !server_id) {
                return;
            }
            try {
                const messageDataArray = await fetchServerMessages(server_id, number_of_messages);
                const newmap = messageDataArray.reduce((map, obj) => {
                    const { channel_id } = obj;
                    if (!map.has(channel_id)) {
                        map.set(channel_id, []);
                    }
                    map.get(channel_id).push(obj);
                    return map;
                }, new Map());

                setChannelMessages(() => newmap);
            } catch (error) {
                console.error("Error fetching messages:", error);
                setChannelMessages(new Map());
            }
        })();
    }, [server_id, number_of_messages]);

    useEffect(() => {
        (async () => {
            const channels = await fetchChannels(server_id);
            setChannels(channels);
        })();
    }, [server_id]);
    useEffect(() => {
        const inChannels = channels.find((channel) => channel.ChannelId === channelId);
        try {
            if (channels.length > 0 && !inChannels) {
                navigate(`/servers/${server_id}/channels/${channels[0].ChannelId}`);
            }
        } catch (error) {
            console.error("Error fetching channels:", error);
            setChannels([]);
        }
    }, [channelId, channels]);

    const ws = useRef<WebSocket | null>(null);
    const onSubmit = (t: SyntheticEvent, inputValue: string): string => {
        t.preventDefault();
        if (inputValue.length === 0) {
            return inputValue;
        } else if (inputValue.length >= 20) {
            return inputValue;
        }
        const stringified = JSON.stringify({
            channel_id: channelId,
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
                    server_id: json.serverid,
                    message: json.message,
                    date: new Date(json.date),
                    author: json.username ?? "User " + json.userid,
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
                        messages.set(channel_id, channel_messages.slice(-number_of_messages));
                    } else {
                        messages.set(channel_id, channel_messages);
                    }
                    return new Map(messages);
                });

                if (newMessage.author_id != auth.authState.user?.id) {
                    toast(`New message from ${newMessage.author}`, {
                        description: newMessage.message,
                        action: {
                            label: "View",
                            onClick: () => {
                                // add navigation to server 2 here
                                navigate(`/servers/2`);
                            },
                        },
                    });
                }
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

    const onChannelSelect = (newChannelId: number) => {
        if (newChannelId === channelId) {
            return;
        }
        navigate(`/servers/${server_id}/channels/${newChannelId}`);
    };

    return (
        <div className="flex flex-grow">
            <div className="flex flex-shrink-0">
                <ChannelSidebar
                    channels={channels}
                    selectedChannelId={channelId}
                    onChannelSelect={onChannelSelect}
                />
            </div>
            <div className="flex flex-grow">
                <ChatPage
                    channel_id={channelId}
                    onSubmit={onSubmit}
                    messages={channnelMessages.get(channelId) || []}
                />
            </div>
        </div>
    );
}

export default ServerPage;
