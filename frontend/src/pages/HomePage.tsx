import { Toaster } from "sonner";
import { toast } from "sonner";
import { useMessageStore } from "@/store/message_store";
import { MessageData } from "@/components/Message";
import { useChannelStore } from "@/store/channel_store";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { useEffect, useState } from "react";
import { useAuth } from "@/AuthContext";
import { serverApi } from "@/api";
import { useNavigate, useParams } from "react-router-dom";
import ServerPage from "./ServerPage";
import { useServerStore } from "@/store/server_store";
import { WebSocketProvider } from "@/WebsocketContext";


export function HomePage() {
  const auth = useAuth();
  const navigate = useNavigate();
  const { serverId: serverIdStr } = useParams<{ serverId: string }>();
  const serverId = serverIdStr ? parseInt(serverIdStr) : -1;
  const [currentServerId, setCurrentServerId] = useState<number | null>(null);
  const servers = useServerStore((state) => state.servers);
  const setServers = useServerStore((state) => state.setServers);
  const [currentName, setCurrentName] = useState<string>("");

  const addChannelMessage = useMessageStore((state) => state.addMessage);

  useEffect(() => {
    auth.addLogoutCallback(() => {
      console.log("logout callback");
    });
  }, []);

  const fetchServers = async () => {
    if (auth.authState.user === null) {
      setCurrentName("");
      setCurrentServerId(null);
      setServers([]);
      return;
    }

    try {
      const retrievedServers = await serverApi.fetchUserServers(auth.authState.user.id);
      // Set the selected server if serverId is provided
      setServers(retrievedServers);
      if (serverId === -1 && retrievedServers.length > 0) {
        navigate(`/servers/${retrievedServers[0].ServerId}`, { replace: true });
      }
    } catch (error) {
      console.error("Error fetching servers:", error);
    }
  };

  useEffect(() => {
    fetchServers();
  }, [auth.authState.user]);

  useEffect(() => {
    if (servers.length === 0) {
      return;
    }
    if (serverId === undefined) {
      navigate(`/servers/${servers[0]?.ServerId}`, { replace: true });
      return;
    }
    if (currentServerId === serverId) {
      return;
    }
    const server = servers.find((s) => s.ServerId === serverId);

    if (!server) {
      return;
    }
    setCurrentServerId(server.ServerId);
    setCurrentName(server?.ServerName || "");
  }, [serverId, servers, navigate]);

  const onMessage = (event: MessageEvent) => {

    const json = JSON.parse(event.data);
    console.log("ws message", json);

    if (json.message_type !== "message") {
      return;
    }
    try {
      const payload = json.payload;
      const newMessage: MessageData = {
        message_id: payload.messageid,
        channel_id: payload.channelid,
        server_id: payload.serverid,
        message: payload.message,
        date: new Date(payload.date),
        author: payload.username ?? "User " + payload.userid,
        author_id: payload.userid,
      };
      const channel_id = newMessage.channel_id;
      if (!channel_id) {
        return;
      }
      addChannelMessage(channel_id, newMessage);
      //if (channelIdRef.current === channel_id) {
      //  return;
      //}
      const currentChannels = useChannelStore.getState().channels;
      const channelName = currentChannels.find((channel) => channel.ChannelId === channel_id)?.ChannelName || "Unknown Channel";
      const max_message_length = 80;
      const shortened_message = newMessage.message.length > max_message_length ? newMessage.message.slice(0, max_message_length) + "..." : newMessage.message;

      if (newMessage.author_id != auth.authState.user?.id) {
        toast(`New message from ${newMessage.author} in ${channelName}`, {
          description: shortened_message,
          action: {
            label: "View",
            onClick: () => {
              navigate(`/servers/${newMessage.server_id}/channels/${newMessage.channel_id}`);
            },
          },
        });
      }
    } catch (err) {
      console.log(err);
    }
  };

  return (
    <WebSocketProvider url={import.meta.env.VITE_WEBSOCKET_URL} onMessage={onMessage}>
      <div className="flex h-full w-full">
        <Toaster />
        <SidebarProvider>
          <AppSidebar serverId={currentServerId} onServerCreate={fetchServers} />
          <div className="flex flex-grow flex-col">
            <div className="sticky top-0 z-10 flex bg-background px-2 py-4 shadow-sm">
              <div className="flex items-center justify-between">
                <SidebarTrigger className="ml-1" />
                <h1 className="text-2xl font-bold">{currentName || ""}</h1>
              </div>
            </div>
            <div className="flex flex-grow flex-col overflow-hidden">
              <div className="flex h-full w-full">
                <ServerPage server_id={serverId} number_of_messages={25} />
              </div>
            </div>
          </div>
        </SidebarProvider>
      </div>
    </WebSocketProvider>
  );
}
