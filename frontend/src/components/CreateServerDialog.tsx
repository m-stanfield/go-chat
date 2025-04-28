import { ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuTrigger } from "@radix-ui/react-context-menu";
import { CreateChannelDialog } from "./CreateChannelDialog";
import { useState } from "react";
import { Button } from "./ui/button";

interface CreateServerDialogProps {
    children: React.ReactNode;
    serverid: number;
}


export function CreateServerDialog({ children, serverid }: CreateServerDialogProps) {

    const [openCreateChannelDialog, setOpenCreateChannelDialog] = useState(false);

    return (

        <>
            <CreateChannelDialog serverid={serverid} open={openCreateChannelDialog} setOpen={setOpenCreateChannelDialog} />
            <ContextMenu>
                <ContextMenuTrigger>
                    {children}
                </ContextMenuTrigger>
                <ContextMenuContent>
                    <ContextMenuItem onSelect={() => setOpenCreateChannelDialog(true)}>
                        Create Channel
                    </ContextMenuItem>
                </ContextMenuContent>
            </ContextMenu>
        </>
    )
}

