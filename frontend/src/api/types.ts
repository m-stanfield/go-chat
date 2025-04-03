// Define all API types in one place
export interface LoginPayload {
  username: string;
  password: string;
}

export interface LoginResponse {
  userid: number;
  // Add other response fields if needed
}

export interface RawMessageData {
  messageid: number;
  channelid: number;
  userid: number;
  username?: string;
  date: string;
  message: string;
}

export interface RawChannelData {
  ChannelId: number;
  ChannelName: string;
}

export interface ServerIconResponse {
  ServerId: number;
  ServerName: string;
  ImageUrl: string | undefined;
} 