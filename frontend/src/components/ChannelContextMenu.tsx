
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

    console.log("ChannelSidebarContextMenu", channelId);

    return (
        <div className={className}>
            <ContextMenu>
                <ContextMenuTrigger>
                    {children}
                </ContextMenuTrigger>
                <ContextMenuContent>
                    <ContextMenuItem>Profile</ContextMenuItem>
                    <ContextMenuItem>Billing</ContextMenuItem>
                    <ContextMenuItem>Team</ContextMenuItem>
                    <ContextMenuItem>Subscription</ContextMenuItem>
                </ContextMenuContent>
            </ContextMenu>
        </div>
    );
}
