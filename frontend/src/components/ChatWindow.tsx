import MessageSubmitWindow from "./MessageSubmitWindow";
import { SyntheticEvent, useEffect, useRef } from "react";
import Message, { MessageData } from "./Message"; // Import your Message component

interface ChatPageProps {
    channel_id: number | undefined;
    messages: MessageData[];
    onSubmit: (t: SyntheticEvent, inputValue: string) => string;
}

function ChatPage({ channel_id, messages, onSubmit }: ChatPageProps) {
    const messageEndRef = useRef<HTMLDivElement>(null);

    // Scroll to bottom whenever messages change or channel changes
    useEffect(() => {
        setTimeout(() => {
            messageEndRef.current?.scrollIntoView({
                behavior: "smooth", // Add smooth scrolling
                block: "end", // Align to the bottom
                inline: "nearest", // Avoid horizontal scrolling
            });
        }, 0);
    }, [messages, channel_id]);

    return (
        <div className="flex flex-grow flex-col rounded-lg bg-gray-600 p-2">
            <div className="flex min-h-0 flex-col overflow-y-auto">
                <ul className="space-y-1 rounded-lg">
                    {messages.map((m) => (
                        <div
                            key={m.message_id}
                            className="flex flex-grow rounded-lg bg-slate-700 hover:bg-slate-600"
                        >
                            <Message message={m} />
                        </div>
                    ))}
                </ul>
                <div ref={messageEndRef} /> {/* Add scroll anchor at the bottom */}
            </div>
            {/* Fixed input at bottom */}
            <div className="mt-2 flex flex-shrink">
                <MessageSubmitWindow onSubmit={onSubmit} />
            </div>
        </div>
    );
}

export default ChatPage;
