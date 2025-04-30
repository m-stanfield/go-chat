

import { useState } from "react"
import {
    Dialog,
} from "@/components/ui/dialog"
import {
    ContextMenu,
    ContextMenuContent,
    ContextMenuItem,
    ContextMenuTrigger,
} from "@/components/ui/context-menu"
import { DialogTrigger } from "@radix-ui/react-dialog";
import { CreateChannelDialog } from "./CreateChannelDialog";

interface SidebarContextMenuProps {
    children: React.ReactNode;
    serverid: number;
    className?: string;
};
export default function SidebarContextMenu({ serverid, children, className }: SidebarContextMenuProps) {

    const [dialogOpen, setDialogOpen] = useState(false)
    return (
        <>

            <Dialog >
                <ContextMenu>
                    <ContextMenuTrigger asChild>
                        <div className={className}>
                            {children}
                        </div >
                    </ContextMenuTrigger>
                    <ContextMenuContent>
                        <DialogTrigger asChild>
                            <ContextMenuItem onClick={() => setDialogOpen(true)}>
                                Create Channel
                            </ContextMenuItem>
                        </DialogTrigger>
                    </ContextMenuContent>
                </ContextMenu>

                <CreateChannelDialog open={dialogOpen} setOpen={setDialogOpen} serverid={serverid} />

            </Dialog>
        </>
    );
}
