import React from "react"; // Import React (assuming you're using React)
import Message, { MessageData } from "./Message"; // Import your Message component

interface MessageViewProps {
    messages: MessageData[];
}
// <div className="h-full w-full flex flex-row  items-end justify-between  bg-slate-800">
export default function MessageView({ messages }: MessageViewProps) {
    const items = messages.map((m) => (
        <li
            key={m.message_id}
            className="w-full rounded-lg bg-slate-700 hover:bg-slate-600"
        >
            <Message message={m} />
        </li>
    ));
    return <ul className="w-full space-y-1  rounded-lg bg-slate-800">{items}</ul>;
}
// items-end
