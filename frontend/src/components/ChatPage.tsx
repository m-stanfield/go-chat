import MessageSubmitWindow from "./MessageSubmitWindow";
import { SyntheticEvent, useEffect, useRef } from "react";
import Message, { MessageData } from "./Message"; // Import your Message component
import { useAuth } from "../AuthContext";

interface ChatPageProps {
    channel_id: number | undefined;
    messages: MessageData[];
    onSubmit: (t: SyntheticEvent, inputValue: string) => string;
}
function ChatPage({ channel_id, messages, onSubmit }: ChatPageProps) {
    const auth = useAuth();
    const messageEndRef = useRef(null);

    useEffect(() => {
        const lastMessage = messages[messages.length - 1];
        if (lastMessage && lastMessage.author_id == auth.authState.user?.id) {
            messageEndRef.current?.scrollIntoView({
                behavior: "smooth",
            });
        }
    }, [auth.authState.user?.id, messages]);

    return (
        <div className="flex flex-grow flex-col overflow-y-auto">
            <div className="">Channel ID: {channel_id}</div>
            <div className="overflow-y-auto">
                <ul className="">
                    {messages.map((m, index) => (
                        <li
                            key={m.message_id}
                            className=""
                            ref={index === messages.length - 1 ? messageEndRef : null}
                        >
                            <Message message={m} />
                        </li>
                    ))}
                </ul>
            </div>
            <div className="">
                <MessageSubmitWindow onSubmit={onSubmit} />
            </div>
        </div>
    );
}

export default ChatPage;
