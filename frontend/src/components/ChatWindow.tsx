import MessageSubmitWindow from "./MessageSubmitWindow";
import { SyntheticEvent, useEffect, useRef } from "react";
import Message from "./Message";
import { useMessageStore } from "@/store/message_store";
import { useWebSocket } from "@/WebsocketContext";

interface ChatPageProps {
    channel_id: number;
}

function ChatPage({ channel_id }: ChatPageProps) {
    const messageEndRef = useRef<HTMLDivElement>(null);
    const messages = useMessageStore((state) => state.messagesByChannel[channel_id]);


    const ws = useWebSocket();
    const onSubmit = (t: SyntheticEvent, inputValue: string): string => {
        t.preventDefault();
        if (inputValue.length === 0) {
            return inputValue;
        } else if (inputValue.length >= 1000) {
            return inputValue;
        }
        const payload = {
            channel_id: channel_id,
            message: inputValue,
        };
        if (ws === null) {
            console.log("websocket hasn't be initialized yet");
            return inputValue;
        }
        ws.sendMessage({ message_type: "channel_message", payload: payload });
        return "";
    };

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


    const validateMessage = (x: string) => {
        if (x.length > 1000) {
            return "Message is too long";
        }
        return "";
    }
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
                    validateMessage={validateMessage}
                />
            </div>
        </div>
    );
}

export default ChatPage;
