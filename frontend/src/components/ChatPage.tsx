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
    const messageEndRef = useRef<HTMLDivElement>(null);

    // Scroll to bottom whenever messages change or channel changes
    useEffect(() => {
        messageEndRef.current?.scrollIntoView();
    }, [messages, channel_id]);

    return (
        <div className="flex flex-col h-full bg-gray-600 p-2 rounded-lg">
            <div className="mb-2">Channel ID: {channel_id}</div>
            {/* Message list with fixed height and scrolling */}
            <div className="flex-grow overflow-y-auto min-h-0">
                <ul className="space-y-1 rounded-lg">
                    {messages.map((m) => (
                        <li
                            key={m.message_id}
                            className="flex flex-grow bg-slate-700 hover:bg-slate-600 rounded-lg"
                        >
                            <Message message={m} />
                        </li>
                    ))}
                </ul>
                <div ref={messageEndRef} /> {/* Add scroll anchor at the bottom */}
            </div>
            {/* Fixed input at bottom */}
            <div className="mt-2">
                <MessageSubmitWindow onSubmit={onSubmit} />
            </div>
        </div>
    );
}

export default ChatPage;
