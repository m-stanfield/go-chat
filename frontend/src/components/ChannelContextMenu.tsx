
import {
    ContextMenu,
    ContextMenuContent,
    ContextMenuItem,
    ContextMenuTrigger,
} from "@/components/ui/context-menu"

interface ChannelSidebarContextMenuProps {
    channelId: number;
    children: React.ReactNode;
};
export default function ChannelSidebarContextMenu({ channelId, children }: ChannelSidebarContextMenuProps) {

    console.log("ChannelSidebarContextMenu", channelId);

    return (
        <>
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
        </>
    );
}
