
import type React from "react"

import { Button } from "@/components/ui/button"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { CreateChannelResponse, postChannel } from "@/api/serverApi"

import { useState } from "react"
import {
    ContextMenuContent,
    ContextMenu,
    ContextMenuItem,
    ContextMenuTrigger,
} from "@/components/ui/context-menu"

interface SidebarContextMenuProps {
    children: React.ReactNode;
    serverid: number;
    className?: string;
    onChannelCreated?: (Channel: CreateChannelResponse) => void;
};
export default function SidebarContextMenu({ serverid, children, className, onChannelCreated }: SidebarContextMenuProps) {

    const [dialogOpen, setDialogOpen] = useState(false)
    const [channelName, setChannelName] = useState("")
    const [isLoading, setIsLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        setIsLoading(true)
        setError(null)
        try {
            // Perform API call to create channel
            const response = await postChannel(serverid, channelName);
            // Reset form and close dialog on success
            setDialogOpen(false)
            onChannelCreated?.(response)
        } catch (err) {
            setError(err instanceof Error ? err.message : "An unknown error occurred")
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
            <ContextMenu modal={false}>
                <ContextMenuTrigger >
                    <div className={className}>
                        {children}
                    </div>
                </ContextMenuTrigger>
                <ContextMenuContent>
                    <ContextMenuItem onSelect={() => setDialogOpen(true)}>
                        Create Channel
                    </ContextMenuItem>
                </ContextMenuContent>
            </ContextMenu>

            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Are you absolutely sure?</DialogTitle>
                    <DialogDescription>
                        This action cannot be undone. Are you sure you want to permanently
                        delete this file?
                    </DialogDescription>
                </DialogHeader>

                <div className="grid gap-4 py-4">
                    <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="name" className="text-right">
                            Name
                        </Label>
                        <Input
                            id="name"
                            value={channelName}
                            onChange={(e) => setChannelName(e.target.value)}
                            className="col-span-3"
                            placeholder="general"
                            required
                        />
                    </div>
                    {error && <div className="text-sm text-red-500 col-span-4 text-center">{error}</div>}
                </div>
                <DialogFooter>
                    <Button type="submit" onClick={handleSubmit} disabled={isLoading}>
                        {isLoading ? "Creating..." : "Create Channel"}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )


}
