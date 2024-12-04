"use client";

import MessageView from "./MessageView";
import MessageSubmitWindow from "./MessageSubmitWindow";
import { MessageData } from "./Message";
import { SyntheticEvent, useEffect, useRef, useState } from "react";

import MockMessages from "./MockMessages";

function ChatPage() {
    const [messages, setMessages] = useState(MockMessages);
    const ws = useRef<WebSocket | null>(null);
    const onmessage = function(event: MessageEvent) {
        const json = JSON.parse(event.data);
        try {
            const newMessage = new MessageData({});
            newMessage.id = json.messageid;
            newMessage.message = json.message;
            newMessage.date = new Date(json.date);
            newMessage.author = json.username;
            setMessages((messages) => [...messages, newMessage]);
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

    useEffect(() => {
        ws.current = new WebSocket("ws://localhost:8080/websocket");
        ws.current.onopen = () => {
            console.log("opening ws");
        };
        ws.current.onmessage = onmessage;
        ws.current.onclose = () => console.log("ws closed");

        return () => {
            ws.current?.close();
        };
    }, []);

    const dummy = useRef<HTMLDivElement | null>(null);
    const [focusMessageWindow, setFocuseMessageWindow] = useState(false);

    const onSubmit = (t: SyntheticEvent, inputValue: string) => {
        t.preventDefault();
        if (inputValue.length === 0) {
            return false;
        }
        setFocuseMessageWindow(true);
        const goMessage = class {
            username: string;
            message: string;
            constructor(UserName: string, Message: string) {
                this.username = UserName;
                this.message = Message;
            }
        };
        const stringified = JSON.stringify({
            ...new goMessage("me", inputValue),
        });
        if (ws?.current === null || ws?.current.readyState === WebSocket.CLOSED) {
            console.log("can't send ws closed");
            return false;
        }
        ws?.current.send(stringified);
        return true;
    };
    return (
        <>
            <div className=" flex h-full w-full flex-col ">
                <div className="flex w-full flex-grow flex-col-reverse overflow-y-scroll rounded-lg">
                    <div ref={dummy}></div>
                    <MessageView messages={messages} />
                </div>
                <div className=" bottom-0 flex w-full p-3  ">
                    <MessageSubmitWindow onSubmit={onSubmit} />
                </div>
            </div>
        </>
    );
}

export default ChatPage;
