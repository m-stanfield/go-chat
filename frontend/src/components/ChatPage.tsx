import MessageSubmitWindow from "./MessageSubmitWindow";
import { SyntheticEvent, useEffect, useRef, useState } from "react";
import Message, { MessageData } from "./Message"; // Import your Message component
import { useAuth } from "../AuthContext";

interface ChatPageProps {
    channel_id: number;
    messages: MessageData[];
    onSubmit: (t: SyntheticEvent, inputValue: string) => string;
}
function ChatPage({ channel_id, messages, onSubmit }: ChatPageProps) {
    const auth = useAuth();
    const messageEndRef = useRef(null);

    useEffect(() => {
        const lastMessage = messages[messages.length - 1];
        if (lastMessage.author_id == auth.authState.user?.id) {
            messageEndRef.current?.scrollIntoView({
                behavior: "smooth",
            });
        }
    }, [auth.authState.user?.id, messages]);
    return (
        <div className=" flex h-full w-full flex-col ">
            <div className="w-full">Channel ID: {channel_id}</div>
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
