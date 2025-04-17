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
import { ChevronUp, User2 } from "lucide-react";
import { DropdownMenuContent } from "./ui/dropdown-menu";
import { postServers } from "@/api/serverApi";
import { CreateServerDialog } from "./CreateServerDialog";

interface SidebarMenuItem {
  title: string;
  url: string;
  icon: React.FC;
  selected: boolean;
}
interface SidebarMenuItemProps {
  items: SidebarMenuItem[];
  onServerCreate: () => void;
}

export function AppSidebar({ items, onServerCreate }: SidebarMenuItemProps) {
  const navigate = useNavigate();
  const auth = useAuth();
  return (
    <Sidebar>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Application</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
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
