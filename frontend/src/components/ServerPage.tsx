import ChatPage from "./ChatPage";
import { SyntheticEvent, useEffect, useRef, useState } from "react";
import MockMessages from "./MockMessages";
import { MessageData } from "./Message";

interface ServerPageProps {
    server_id: number;
}
function ServerPage({ server_id }: ServerPageProps) {
    const selectedChannelId = 2;
    const [messages, setMessages] = useState<MessageData[]>(MockMessages);
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
                    id: json.messageid,
                    channel_id: json.channel_id,
                    message: json.message,
                    date: new Date(json.date),
                    author: json.username,
                    author_id: json.userid,
                };
                setMessages((messages) => {
                    console.log(messages.length);
                    if (messages.length > maxMessageLength) {
                        return [...messages.slice(-maxMessageLength), newMessage];
                    }
                    return [...messages, newMessage];
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
        <div className="w-full h-full">
            <div>
                <h1 className="w-full">Server ID: {server_id}</h1>
            </div>
            <div className="h-full w-full">
                <ChatPage
                    channel_id={selectedChannelId}
                    onSubmit={onSubmit}
                    messages={messages}
                />
            </div>
        </div>
    );
}

export default ServerPage;
