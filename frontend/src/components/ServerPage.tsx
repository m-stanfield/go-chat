import ChatPage from "./ChatPage";
import { SyntheticEvent, useEffect, useRef, useState } from "react";
import { MessageData } from "./Message";
import IconBanner, { IconInfo } from "./IconList";
import { fetchServerMessages, fetchChannels } from "../api/serverApi";
import { useParams } from "react-router-dom";

interface ServerPageProps {
    number_of_messages: number;
}

function ServerPage({ number_of_messages }: ServerPageProps) {
    const { serverId } = useParams();
    const server_id = parseInt(serverId || "-1");
    const [selectedChannelId, setSelectedChannelId] = useState<number>(0);
    const [channnelMessages, setChannelMessages] = useState<Map<number, MessageData[]>>(new Map());
    const [channelInformationArray, setChannelInformationArray] = useState<IconInfo[]>([]);
    const ws = useRef<WebSocket | null>(null);

    // Reset state when server changes
    useEffect(() => {
        setSelectedChannelId(0);
        setChannelMessages(new Map());
        setChannelInformationArray([]);
        
        // Fetch channels for new server
        const fetchChannelsForServer = async () => {
            if (server_id < 0) {
                return;
            }
            try {
                const channelInfoArray = await fetchChannels(server_id);
                setChannelInformationArray(channelInfoArray);
                
                if (channelInfoArray.length > 0) {
                    setSelectedChannelId(channelInfoArray[0].icon_id);
                }
            } catch (error) {
                console.error("Error fetching channels:", error);
                setChannelInformationArray([]);
            }
        };

        fetchChannelsForServer();
    }, [server_id]); // Only re-run when server_id changes

    // Fetch messages when server or channel changes
    useEffect(() => {
        const fetchMessages = async () => {
            if (server_id < 0) {
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
                
                setChannelMessages(newmap);
            } catch (error) {
                console.error("Error fetching messages:", error);
                setChannelMessages(new Map());
            }
        };

        fetchMessages();
    }, [server_id, selectedChannelId, number_of_messages]);

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

    return (
        <div className="flex flex-row h-full">
            <div className="mr-4 h-full">
                <div className="sticky top-0 p-2">
                    <IconBanner
                        icon_info={channelInformationArray}
                        onServerSelect={setSelectedChannelId}
                        direction="vertical"
                        displayMode="text"
                        selectedIconId={selectedChannelId}
                    />
                </div>
            </div>
            <div className="flex-grow">
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
