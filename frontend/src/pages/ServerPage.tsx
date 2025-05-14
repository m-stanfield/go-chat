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
import { Channel } from "@/types/channel";
import { useAuth } from "@/AuthContext";
import { toast } from "sonner";
import { useNavigate, useParams } from "react-router-dom";
import SidebarContextMenu from "@/components/SidebarContextMenu";
import { useMessageStore } from "@/store/message_store";

interface ServerPageProps {
  server_id: number;
  number_of_messages: number;
}

function ServerPage({ server_id, number_of_messages }: ServerPageProps) {
  const auth = useAuth();
  const navigate = useNavigate();
  const { channelId: channelIdStr } = useParams<{ channelId: string }>();
  const channelId = channelIdStr ? parseInt(channelIdStr) : -1;

  const [channels, setChannels] = useState<Channel[]>([]);
  // get channel messgages from store
  const channelMessages = useMessageStore((state) => state.messagesByChannel);
  const setChannelMessages = useMessageStore((state) => state.setMessagesByChannel);
  const addChannelMessage = useMessageStore((state) => state.addMessage);
  const removeAllMessages = useMessageStore((state) => state.removeAllMessages);
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
      const channels = await fetchChannels(server_id);
      setChannels(channels);
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

  const ws = useRef<WebSocket | null>(null);
  const onSubmit = (t: SyntheticEvent, inputValue: string): string => {
    t.preventDefault();
    if (inputValue.length === 0) {
      return inputValue;
    } else if (inputValue.length >= 1000) {
      return inputValue;
    }
    const stringified = JSON.stringify({
      channel_id: channelId,
      message: inputValue,
    });
    if (ws === null) {
      console.log("websocket hasn't be initialized yet");
      return inputValue;
    }
    if (ws.current?.readyState === WebSocket.CLOSED) {
      console.log("can't send ws closed");
      return inputValue;
    }
    ws.current?.send(stringified);
    return "";
  };

  useEffect(() => {
    const startTime = Date.now();
    console.log("starting to open websocket", Date.now() - startTime);
    ws.current = new WebSocket(`ws://localhost:8080/websocket`);
    ws.current.onopen = () => {
      console.log("opening ws", Date.now() - startTime);
    };
    ws.current.onmessage = function(event: MessageEvent) {
      console.log("ws message", Date.now() - startTime);
      const json = JSON.parse(event.data);
      try {
        const newMessage: MessageData = {
          message_id: json.messageid,
          channel_id: json.channelid,
          server_id: json.serverid,
          message: json.message,
          date: new Date(json.date),
          author: json.username ?? "User " + json.userid,
          author_id: json.userid,
        };
        const channel_id = newMessage.channel_id;
        if (!channel_id) {
          return;
        }
        addChannelMessage(channel_id, newMessage);
        if (newMessage.author_id != auth.authState.user?.id) {
          toast(`New message from ${newMessage.author}`, {
            description: newMessage.message,
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
    ws.current.onclose = () => {
      console.log("ws closed", Date.now() - startTime);
    };

    return () => {
      ws.current?.close();
    };
  }, []);

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
              channels={channels}
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
          onSubmit={onSubmit}
          messages={channelMessages[channelId] || []}
        />
      </div>
    </div>
  );
}

export default ServerPage;
