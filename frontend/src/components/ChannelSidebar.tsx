import ChannelSidebarContextMenu from "./ChannelContextMenu";
import { useChannelStore } from "@/store/channel_store";

interface ChannelSidebarProps {
    selectedChannelId: number;
    serverid: number;
    onChannelSelect: (channelId: number) => void;
}

export default function ChannelSidebar({
    selectedChannelId,
    onChannelSelect,
}: ChannelSidebarProps) {

    const channels = useChannelStore((state) => state.channels);

    return (
        <>
            <div className="flex flex-grow h-full flex-col w-48 bg-gray-800 p-4 text-white">
                <h2 className="mb-4 text-xl font-bold">Channels</h2>
                <div className="flex flex-grow items-start mb-4 overflow-y-auto overflow-x-hidden">
                    <ul className="space-y-2">
                        {channels.map((channel) => (
                            <ChannelSidebarContextMenu key={channel.ChannelId} channelId={channel.ChannelId} >
                                <li
                                    key={channel.ChannelId}
                                    onClick={() => onChannelSelect(channel.ChannelId)}
                                    className={`flex  flex-grow cursor-pointer rounded p-2 hover:bg-gray-600 ${channel.ChannelId === selectedChannelId ? "bg-gray-700" : ""}`}
                                >
                                    # {channel.ChannelName}
                                </li>
                            </ChannelSidebarContextMenu>
                        ))}
                    </ul>
                </div >

            </div >
        </>
    );
}
