
import { create } from 'zustand';

import { MessageData } from "@/components/Message";



type MessageState = {
    messagesByChannel: Record<number, MessageData[]>
    setMessagesByChannel: (channelId: number, messages: MessageData[]) => void
    addMessage: (channelId: number, message: MessageData) => void
    removeMessage: (channelId: number, messageId: number) => void
    removeAllMessages: () => void
}

export const useMessageStore = create<MessageState>((set) => ({
    messagesByChannel: {},

    setMessagesByChannel: (channelId: number, messages: MessageData[]) =>
        set((state) => {
            return {
                messagesByChannel: {
                    ...state.messagesByChannel,
                    [channelId]: messages,
                },
            }
        }),


    addMessage: (channelId, message) =>
        set((state) => {
            const existingMessages = state.messagesByChannel[channelId] || []
            return {
                messagesByChannel: {
                    ...state.messagesByChannel,
                    [channelId]: [...existingMessages, message],
                },
            }
        }),

    removeMessage: (channelId, messageId) =>
        set((state) => {
            const existingMessages = state.messagesByChannel[channelId] || []
            return {
                messagesByChannel: {
                    ...state.messagesByChannel,
                    [channelId]: existingMessages.filter((msg) => msg.message_id !== messageId),
                },
            }
        }),
    removeAllMessages: () =>
        set(() => {
            return {
                messagesByChannel: {
                },
            }
        }),
}))
