import MessageSubmitWindow from "./MessageSubmitWindow";
import { SyntheticEvent, useEffect, useRef } from "react";
import Message from "./Message";
import { useMessageStore } from "@/store/message_store";

interface ChatPageProps {
    channel_id: number;
    onSubmit: (t: SyntheticEvent, inputValue: string) => string;
}

function ChatPage({ channel_id, onSubmit }: ChatPageProps) {
    const messageEndRef = useRef<HTMLDivElement>(null);
    const messages = useMessageStore((state) => state.messagesByChannel[channel_id]);


    // Scroll to bottom whenever messages change or channel changes
    useEffect(() => {
        setTimeout(() => {
            messageEndRef.current?.scrollIntoView({});
        }, 0);
    }, [channel_id]);

    useEffect(() => {
        setTimeout(() => {
            messageEndRef.current?.scrollIntoView({
                behavior: "smooth",
                block: "end",
                inline: "nearest",
            });
        }, 0);
    }, [messages]);

    return (
        <div className="flex h-full w-full flex-col rounded-lg bg-gray-600 p-2">
            <div className="flex flex-1 flex-col overflow-y-scroll">
                <div className="flex flex-1 flex-col justify-end">
                    <ul className="space-y-1 p-2">
                        {messages && (messages).map((m) => (
                            <div
                                key={m.message_id}
                                className="flex flex-grow rounded-lg bg-slate-700 hover:bg-slate-600"
                            >
                                <Message message={m} />
                            </div>
                        ))}
                    </ul>
                    <div ref={messageEndRef} />
                </div>
            </div>
            <div className="flex-shrink-0 p-2">
                <MessageSubmitWindow
                    onSubmit={onSubmit}
                    validateMessage={(x) => {
                        if (x.length > 1000) {
                            return "Message is too long";
                        }
                        return "";
                    }}
                />
            </div>
        </div>
    );
}

export default ChatPage;
