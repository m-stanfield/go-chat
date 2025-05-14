import { create } from 'zustand';
import { ServerData } from '@/api';

type ServerState = {
    servers: ServerData[]
    currentServer: number | null
    setCurrentServer: (serverId: number | null) => void
    setServers: (servers: ServerData[]) => void
    addServer: (server: ServerData) => void
    removeServer: (serverId: number) => void
}

export const useServerStore = create<ServerState>((set) => ({
    servers: [],
    currentServer: null,
    setCurrentServer: (serverId) => set({ currentServer: serverId }),

    setServers: (servers) => set({ servers }),

    addServer: (server) =>
        set((state) => ({
            servers: [...state.servers, server],
        })),

    removeServer: (serverId) =>
        set((state) => ({
            servers: state.servers.filter((msg) => msg.ServerId !== serverId),
        })),
}))
