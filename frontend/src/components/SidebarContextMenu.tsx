
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

interface SidebarContextMenuProps {
    children: React.ReactNode;
    serverid: number;
    open: boolean;
    setOpen: (open: boolean) => void;
    className?: string;
    onChannelCreated?: (Channel: CreateChannelResponse) => void;
};
export default function SidebarContextMenu({ serverid, children, className, onChannelCreated, open, setOpen }: SidebarContextMenuProps) {

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
            setOpen(false)
            onChannelCreated?.(response)
        } catch (err) {
            setError(err instanceof Error ? err.message : "An unknown error occurred")
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <div className={className}>
                {children}
            </div>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Create New Channel</DialogTitle>
                    <DialogDescription>
                        Enter a name for your new channel. Click create when you're done.
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
