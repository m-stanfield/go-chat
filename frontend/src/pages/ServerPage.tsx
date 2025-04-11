import { SyntheticEvent, useEffect, useRef, useState } from "react";
import { fetchServerMessages, fetchChannels } from "../api/serverApi";
import ChatPage from "@/components/ChatWindow";
import { MessageData } from "@/components/Message";
import ChannelSidebar from "@/components/ChannelSidebar";
import { Channel } from "@/types/channel";

interface ServerPageProps {
    server_id: number;
    number_of_messages: number;
}

function ServerPage({ server_id, number_of_messages }: ServerPageProps) {
    const [selectedChannelId, setSelectedChannelId] = useState<number>(1);
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

    useEffect(() => {
        async function getChannels() {
            if (server_id < 0 || !server_id) return;

            try {
                const channels = await fetchChannels(server_id);
                setChannels(channels);
                if (channels.length > 0 && !selectedChannelId) {
                    setSelectedChannelId(channels[0].ChannelId);
                }
            } catch (error) {
                console.error("Error fetching channels:", error);
                setChannels([]);
            }
        }

        getChannels();
    }, [server_id]);

    return (
        <div className="flex flex-1">
            <ChannelSidebar
                channels={channels}
                selectedChannelId={selectedChannelId}
                onChannelSelect={setSelectedChannelId}
            />
            <div className="flex flex-1">
                <ChatPage
                    channel_id={selectedChannelId}
                    onSubmit={onSubmit}
                    messages={channnelMessages.get(selectedChannelId) || []}
                />
            </div>
        </div>
    );
}

export default ServerPage;
