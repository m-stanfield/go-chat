
import {
    ContextMenu,
    ContextMenuContent,
    ContextMenuItem,
    ContextMenuTrigger,
} from "@/components/ui/context-menu"

interface ChannelSidebarContextMenuProps {
    channelId: number;
    children: React.ReactNode;
    className?: string;
};
export default function ChannelSidebarContextMenu({ channelId, children, className }: ChannelSidebarContextMenuProps) {


    return (
        <div className={className}>
            <ContextMenu>
                <ContextMenuTrigger>
                    {children}
                </ContextMenuTrigger>
                <ContextMenuContent>
                    <ContextMenuItem>Channel ID: {channelId}</ContextMenuItem>
                </ContextMenuContent>
            </ContextMenu>
        </div>
    );
}
