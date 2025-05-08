
import { create } from 'zustand';

import { MessageData } from "@/components/Message";


type MessageState = {
    messagesByChannel: Record<string, MessageData[]>
    addMessage: (channelId: number, message: MessageData) => void
    removeMessage: (channelId: number, messageId: number) => void
}

export const useMessageStore = create<MessageState>((set) => ({
    messagesByChannel: {},


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
}))
