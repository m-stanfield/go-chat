import { Channel } from '@/types/channel';
import { create } from 'zustand';

type ChannelState = {
    channels: Channel[]
    setChannels: (channels: Channel[]) => void
    addChannel: (channel: Channel) => void
    removeChannel: (channelId: number) => void
}

export const useChannelStore = create<ChannelState>((set) => ({
    channels: [],

    setChannels: (channels) => set({ channels }),

    addChannel: (channel) =>
        set((state) => ({
            channels: [...state.channels, channel],
        })),

    removeChannel: (channelId) =>
        set((state) => ({
            channels: state.channels.filter((msg) => msg.ChannelId !== channelId),
        })),
}))
