import { Channel } from "@/types/channel";
import { ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuTrigger } from "@radix-ui/react-context-menu";
import { CreateServerDialog } from "./CreateServerDialog";
import { useState } from "react";
import { CreateChannelDialog } from "./CreateChannelDialog";
import { Button } from "./ui/button";

interface ChannelSidebarProps {
    channels: Channel[];
    selectedChannelId: number;
    serverid: number;
    onChannelSelect: (channelId: number) => void;
}

export default function ChannelSidebar({
    channels,
    serverid,
    selectedChannelId,
    onChannelSelect,
}: ChannelSidebarProps) {

    const [openCreateChannelDialog, setOpenCreateChannelDialog] = useState(false);
    const setOpen = (open: boolean) => {
        console.log("setOpen", open);
        setOpenCreateChannelDialog(open);
    }
    return (

        <>
            <CreateChannelDialog serverid={serverid} open={openCreateChannelDialog} setOpen={setOpen} />
            <CreateServerDialog serverid={serverid}>
                <div className="flex flex-grow flex-col w-48 bg-gray-800 p-4 text-white">
                    <h2 className="mb-4 text-xl font-bold">Channels</h2>
                    <div className="flex flex-grow items-start mb-4 overflow-auto">

                        <ContextMenu>
                            <ContextMenuTrigger>
                                <ul className="space-y-2">
                                    {channels.map((channel) => (
                                        <li
                                            key={channel.ChannelId}
                                            className={`cursor-pointer rounded p-2 hover:bg-gray-700 ${channel.ChannelId === selectedChannelId ? "bg-gray-700" : ""
                                                }`}
                                            onClick={() => onChannelSelect(channel.ChannelId)}
                                        >
                                            # {channel.ChannelName}
                                        </li>
                                    ))}
                                </ul>
                            </ContextMenuTrigger>
                            <ContextMenuContent className="bg-gray-800 text-white">
                                <ContextMenuItem onSelect={() => setOpen(true)}>
                                    Create Channel
                                </ContextMenuItem>
                            </ContextMenuContent>
                        </ContextMenu>
                    </div>
                    <div className="flex flex-col align-bottom mt-4">
                    </div>
                </div>
            </CreateServerDialog>
        </>
    );
}
