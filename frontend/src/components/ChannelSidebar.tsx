import { Channel } from "@/types/channel";

interface ChannelSidebarProps {
    channels: Channel[];
    selectedChannelId: number;
    onChannelSelect: (channelId: number) => void;
}

export default function ChannelSidebar({
    channels,
    selectedChannelId,
    onChannelSelect,
}: ChannelSidebarProps) {
    return (
        <div className="w-64 bg-gray-800 p-4 text-white">
            <h2 className="mb-4 text-xl font-bold">Channels</h2>
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
        </div>
    );
}
