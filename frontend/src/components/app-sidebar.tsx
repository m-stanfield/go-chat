import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { Button } from "./ui/button";
import { toast } from "sonner";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/AuthContext";

interface SidebarMenuItem {
  title: string;
  url: string;
  icon: React.FC;
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
                    <div>
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
                          onClick: () => console.log("Undo"),
                        },
                      })
                    }
                  >
                    Toast
                  </Button>
                </SidebarMenuButton>
              </SidebarMenuItem>

              <SidebarMenuItem key="logout-button">
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
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  );
}
