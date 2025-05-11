import { create } from 'zustand';
import { ServerData } from '@/api';

type ServerState = {
    servers: ServerData[]
    setServers: (servers: ServerData[]) => void
    addServer: (server: ServerData) => void
    removeServer: (serverId: number) => void
}

export const useServerStore = create<ServerState>((set) => ({
    servers: [],

    setServers: (servers) => set({ servers }),

    addServer: (server) =>
        set((state) => ({
            servers: [...state.servers, server],
        })),

    removeServer: (serverId) =>
        set((state) => ({
            servers: state.servers.filter((msg) => msg.ServerId !== serverId), // <- check this key
        })),
}))
