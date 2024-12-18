import MessageSubmitWindow from "./MessageSubmitWindow";
import { SyntheticEvent, useEffect, useRef, useState } from "react";
import Message, { MessageData } from "./Message"; // Import your Message component

import MockMessages from "./MockMessages";
import { useAuth } from "../AuthContext";

interface ChatPageProps {
    channel_id: number;
}
function ChatPage({ channel_id }: ChatPageProps) {
    const auth = useAuth();
    const [messages, setMessages] = useState<MessageData[]>(MockMessages);
    const ws = useRef<WebSocket | null>(null);
    const dummy = useRef<HTMLDivElement | null>(null);
    const [focusMessageWindow, setFocuseMessageWindow] = useState(false);
    const maxMessageLength = 30;

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
                    channel_id: channel_id,
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

            if (focusMessageWindow) {
                console.log("updating focus message: " + focusMessageWindow);
                setFocuseMessageWindow(false);
                if (dummy.current) {
                    dummy.current.scrollIntoView({ behavior: "smooth" });
                }
            }
        };
        ws.current.onclose = () => {
            console.log("ws closed", Date.now() - startTime);
        };

        return () => {
            ws.current?.close();
        };
    }, []);

    const onSubmit = (t: SyntheticEvent, inputValue: string): string => {
        t.preventDefault();
        if (inputValue.length === 0) {
            return inputValue;
        }
        setFocuseMessageWindow(true);
        const stringified = JSON.stringify({
            channel_id: channel_id,
            message: inputValue,
        });
        if (ws?.current === null) {
            console.log("websocket hasn't be initialized yet");
            return inputValue;
        }
        if (ws.current.readyState === WebSocket.CLOSED) {
            console.log("can't send ws closed");
            return inputValue;
        }
        ws.current.send(stringified);
        return "";
    };

    const messageEndRef = useRef(null);
    useEffect(() => {
        const lastMessage = messages[messages.length - 1];
        if (lastMessage.author_id == auth.authState.user?.id) {
            messageEndRef.current?.scrollIntoView({
                behavior: "smooth",
            });
        }
    }, [messages]);
    return (
        <div className=" flex h-full w-full flex-col ">
            <div className="flex w-full flex-grow flex-col-reverse overflow-y-scroll rounded-lg">
                <ul className="w-full space-y-1  rounded-lg bg-slate-800">
                    {messages.map((m, index) => (
                        <li
                            key={m.id}
                            className="w-full rounded-lg bg-slate-700 hover:bg-slate-600"
                            ref={index == messages.length - 1 ? messageEndRef : null}
                        >
                            <Message message={m} />
                        </li>
                    ))}{" "}
                </ul>
            </div>
            <div className=" bottom-0 flex w-full p-3  ">
                <MessageSubmitWindow onSubmit={onSubmit} />
            </div>
        </div>
    );
}

export default ChatPage;
