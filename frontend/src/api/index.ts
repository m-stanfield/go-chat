// Create a centralized API client
import { LoginPayload, LoginResponse } from "./types";

const API_BASE_URL = "/api";

export interface ServerData {
  ServerId: number;
  ServerName: string;
  ImageUrl: string;
}

// Create a reusable fetch wrapper with common options
const apiFetch = async (endpoint: string, options: RequestInit = {}) => {
  const defaultOptions: RequestInit = {
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
  };

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...defaultOptions,
    ...options,
    headers: {
      ...defaultOptions.headers,
      ...options.headers,
    },
  });

  if (!response.ok) {
    throw new Error(`API error: ${response.statusText}`);
  }

  return response.json();
};

// Auth API
export const authApi = {
  login: async (payload: LoginPayload): Promise<LoginResponse> => {
    return apiFetch("/auth/login", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },

  register: async (payload: LoginPayload): Promise<void> => {
    return apiFetch("/users", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },
};

// Server API
export const serverApi = {
  fetchMessages: async (serverId: number, messageCount: number) => {
    const data = await apiFetch(`/servers/${serverId}/messages?count=${messageCount}`);
    const messageDataArray = data.messages.map((msg: any) => ({
      ...msg,
      username: msg.username ? msg.username : "User" + msg.userid,
    }));
    return messageDataArray.sort((a: any, b: any) => a.messageid - b.messageid);
  },

  fetchChannels: async (serverId: number) => {
    const data = await apiFetch(`/servers/${serverId}/channels`);
    return data.channels;
  },

  fetchUserServers: async (userId: number): Promise<ServerData[]> => {
    const data = await apiFetch(`/users/${userId}/servers`);
    return data.servers.map((server: { ServerId: string; ServerName: string }) => ({
      ServerId: server.ServerId,
      ServerName: server.ServerName,
      ImageUrl: "https://miro.medium.com/v2/resize:fit:720/format:webp/0*UD_CsUBIvEDoVwzc.png",
    }));
  },
};
