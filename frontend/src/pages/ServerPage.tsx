
import { SyntheticEvent, useEffect, useRef, useState } from "react";
import {
  ContextMenuContent,
  ContextMenu,
  ContextMenuItem,
  ContextMenuTrigger,
} from "@/components/ui/context-menu"
import { fetchServerMessages, fetchChannels, CreateChannelResponse } from "../api/serverApi";
import ChatPage from "@/components/ChatWindow";
import { MessageData } from "@/components/Message";
import ChannelSidebar from "@/components/ChannelSidebar";
import { useAuth } from "@/AuthContext";
import { toast } from "sonner";
import { useNavigate, useParams } from "react-router-dom";
import SidebarContextMenu from "@/components/SidebarContextMenu";
import { useMessageStore } from "@/store/message_store";
import { useChannelStore } from "@/store/channel_store";
import { useWebSocket } from "@/WebsocketContext";

interface ServerPageProps {
  server_id: number;
  number_of_messages: number;
}

function ServerPage({ server_id, number_of_messages }: ServerPageProps) {
  const navigate = useNavigate();
  const { channelId: channelIdStr } = useParams<{ channelId: string }>();
  const channelId = channelIdStr ? parseInt(channelIdStr) : -1;

  const channels = useChannelStore((state) => state.channels);
  const setChannels = useChannelStore((state) => state.setChannels);
  // get channel messgages from store
  const setChannelMessages = useMessageStore((state) => state.setMessagesByChannel);
  const removeAllMessages = useMessageStore((state) => state.removeAllMessages);
  const channelIdRef = useRef<number | undefined>(channelId);

  useEffect(() => {
    channelIdRef.current = channelId;
  }, [channelId]);
  useEffect(() => {
    (async () => {
      if (server_id < 0 || !server_id) {
        return;
      }
      try {
        const messageDataArray = await fetchServerMessages(server_id, number_of_messages);
        const newmap = messageDataArray.reduce((map, obj) => {
          const { channel_id } = obj;
          if (!map.has(channel_id)) {
            map.set(channel_id, []);
          }
          map.get(channel_id).push(obj);
          return map;
        }, new Map());

        newmap.forEach((messages, channel_id) => {
          setChannelMessages(channel_id, messages);
        })
      } catch (error) {
        console.error("Error fetching messages:", error);
        removeAllMessages();

      }
    })();
  }, [server_id, number_of_messages]);

  useEffect(() => {
    (async () => {
      const retrieved_channels = await fetchChannels(server_id)
      setChannels(retrieved_channels);
    })();
  }, [server_id]);
  useEffect(() => {
    const inChannels = channels.find((channel) => channel.ChannelId === channelId);
    try {
      if (channels.length > 0 && !inChannels) {
        navigate(`/servers/${server_id}/channels/${channels[0].ChannelId}`);
      }
    } catch (error) {
      console.error("Error fetching channels:", error);
      setChannels([]);
    }
  }, [channelId, channels]);



  const onChannelSelect = (newChannelId: number) => {
    if (newChannelId === channelId) {
      return;
    }
    navigate(`/servers/${server_id}/channels/${newChannelId}`);
  };

  const onChannelCreated = (newChannel: CreateChannelResponse) => {

    onChannelSelect(newChannel.channelId);
  }

  const [open, setOpen] = useState(false)

  return (
    <div className="flex flex-grow">
      <SidebarContextMenu serverid={server_id} className="flex h-full" onChannelCreated={onChannelCreated} open={open} setOpen={setOpen}>

        <ContextMenu modal={false}>
          <ContextMenuTrigger >
            <ChannelSidebar
              serverid={server_id}
              selectedChannelId={channelId}
              onChannelSelect={onChannelSelect}
            />

          </ContextMenuTrigger>
          <ContextMenuContent>
            <ContextMenuItem onSelect={() => setOpen(true)}>
              Create Channel
            </ContextMenuItem>
          </ContextMenuContent>
        </ContextMenu>
      </SidebarContextMenu>
      <div className="flex flex-grow">
        <ChatPage
          channel_id={channelId}
        />
      </div>
    </div>
  );
}

export default ServerPage;
