

import {
    ContextMenu,
    ContextMenuContent,
    ContextMenuItem,
    ContextMenuTrigger,
} from "@/components/ui/context-menu"

interface SidebarContextMenuProps {
    children: React.ReactNode;
    className?: string;
};
export default function SidebarContextMenu({ children, className }: SidebarContextMenuProps) {
    return (
        <ContextMenu>
            <ContextMenuTrigger asChild>
                <div className={className}>
                    {children}
                </div >
            </ContextMenuTrigger>
            <ContextMenuContent>
                <ContextMenuItem>Channel Sidebar</ContextMenuItem>
            </ContextMenuContent>
        </ContextMenu>
    );
}
