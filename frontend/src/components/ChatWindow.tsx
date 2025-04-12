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
            messageEndRef.current?.scrollIntoView({});
        }, 0);
    }, [messages, channel_id]);

    return (
        <div className="flex h-full w-full flex-col rounded-lg bg-gray-600 p-2">
            <div className="flex flex-1  overflow-hidden">
                <div className="flex-1 overflow-y-scroll">
                    <ul className="space-y-1 p-2">
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
            </div>
            {/* Fixed input at bottom */}
            <div className="flex-shrink-0 p-2">
                <MessageSubmitWindow onSubmit={onSubmit} />
            </div>
        </div>
    );
}

export default ChatPage;
