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

interface SidebarMenuItem {
  title: string;
  url: string;
  icon: React.FC;
  selected: boolean;
}
interface SidebarMenuItemProps {
  items: SidebarMenuItem[];
}

export function AppSidebar({ items }: SidebarMenuItemProps) {
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

              <SidebarMenuItem key="toast-button">
                <SidebarMenuButton asChild>
                  <Button
                    onClick={() =>
                      toast("Event has been created", {
                        description: "Sunday, December 03, 2023 at 9:00 AM",
                        action: {
                          label: "Undo",
                          onClick: () => {
                            // generate random number between 0 and 1
                            const randomNumber = Math.random();
                            navigate(`/servers/${randomNumber > 0.5 ? 1 : 2}`);
                          },
                        },
                      })
                    }
                  >
                    Toast
                  </Button>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
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
                  <SidebarMenuButton asChild>
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
