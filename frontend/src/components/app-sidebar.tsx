import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarFooter,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { Button } from "./ui/button";
import { toast } from "sonner";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/AuthContext";
import { DropdownMenu, DropdownMenuItem, DropdownMenuTrigger } from "@radix-ui/react-dropdown-menu";
import { ChevronUp, Home, User2 } from "lucide-react";
import { DropdownMenuContent } from "./ui/dropdown-menu";
import { CreateServerDialog } from "./CreateServerDialog";
import { useEffect, useState } from "react";
import { ServerData } from "@/api";
import { useServerStore } from "@/store/server_store";

interface SidebarMenuItem {
  title: string;
  url: string;
  icon: React.FC;
  selected: boolean;
}
interface SidebarMenuItemProps {
  serverId: number | null;
  onServerCreate: () => void;
}

export function AppSidebar({ serverId, onServerCreate }: SidebarMenuItemProps) {
  const navigate = useNavigate();
  const auth = useAuth();

  const servers = useServerStore((state) => state.servers);

  const [serverNavItem, setServerNavItem] = useState<SidebarMenuItem[]>([]);
  useEffect(() => {
    const serverNavItems = servers.map((server: ServerData) => ({
      title: server.ServerName,
      url: `/servers/${server.ServerId}`,
      icon: Home,
      selected: server.ServerId === serverId,
    }));
    setServerNavItem(serverNavItems);
  }, [servers, serverId]);
  return (
    <Sidebar>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Application</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {serverNavItem.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild onClick={() => navigate(item.url)}>
                    <div
                      className={`${item.selected ? "bg-slate-400 hover:bg-blue-800" : "hover:bg-blue-800"}`}
                    >
                      <item.icon />
                      <span>{item.title}</span>
                    </div>
                  </SidebarMenuButton>

                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuButton asChild
            className="flex align-bottom my-1"
          >
            <CreateServerDialog onServerCreated={onServerCreate} />
          </SidebarMenuButton>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton>
                  <User2 /> @{auth.authState.user?.name}
                  <ChevronUp className="ml-auto" />
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent side="top" className="w-[--radix-popper-anchor-width]">
                <DropdownMenuItem>

                  <SidebarMenuButton asChild
                    className="my-1"
                  >
                    <Button
                      onClick={() => {
                        auth.logout();
                        toast("Logged out successfully", {
                          description: "You have been logged out",
                        });
                        navigate("/login");
                      }}
                    >
                      Logout
                    </Button>
                  </SidebarMenuButton>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  );
}
