import { Toaster } from "sonner";
import { Home } from "lucide-react";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { useEffect, useState } from "react";
import { useAuth } from "@/AuthContext";
import { serverApi, ServerData } from "@/api";
import { useNavigate, useParams } from "react-router-dom";

interface ServerNavItem {
  title: string;
  url: string;
  icon: React.FC;
  selected: boolean;
}

export function HomePage() {
  const auth = useAuth();
  const navigate = useNavigate();
  const { serverId } = useParams<{ serverId: string }>();
  const [serverNaveItem, setServerNavItem] = useState<ServerNavItem[]>([]);
  const [serverdata, setServers] = useState<ServerData[]>([]);
  const [currentName, setCurrentName] = useState<string>("");

  useEffect(() => {
    auth.addLogoutCallback(() => {
      console.log("logout callback");
    });
  }, []);

  useEffect(() => {
    const fetchServers = async () => {
      if (auth.authState.user === null) {
        setCurrentName("");
        setServers([]);
        return;
      }

      try {
        const servers = await serverApi.fetchUserServers(auth.authState.user.id);
        // Set the selected server if serverId is provided
        setServers(servers);
      } catch (error) {
        console.error("Error fetching servers:", error);
      }
    };

    fetchServers();
  }, [auth.authState.user]);

  useEffect(() => {
    const server = serverdata.find((s) => s.ServerId === parseInt(serverId || ""));
    if (serverId) {
      setCurrentName(server?.ServerName || "");
    } else {
      navigate(`/servers/${serverdata[0]?.ServerId}`);
    }
    const serverNavItems = serverdata.map((server: ServerData) => ({
      title: server.ServerName,
      url: `/servers/${server.ServerId}`,
      icon: Home,
      selected: server.ServerId === parseInt(serverId || ""),
    }));
    setServerNavItem(serverNavItems);
  }, [serverId, serverdata]);

  // create placeholder text that is 10k characters
  const placeholderText = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. ".repeat(2000);

  return (
    <div className="flex h-full w-full">
      <Toaster />
      <SidebarProvider>
        <AppSidebar items={serverNaveItem} />
        <div className="flex flex-1 flex-col">
          <div className="sticky top-0 z-10 bg-background px-2 py-4 shadow-sm">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <SidebarTrigger className="ml-1" />
                <h1 className="text-2xl font-bold">{currentName || ""}</h1>
              </div>
            </div>
          </div>
          <div className="rounded-lg bg-card p-6 shadow">
            <h1 className="mb-4 text-2xl font-bold text-card-foreground">Home Page</h1>
            <p className="text-muted-foreground">Welcome to the home page of the application.</p>
            <div className="mt-6 grid gap-4 md:grid-cols-2">
              <div className="rounded-md bg-primary/10 p-4">
                <h3 className="font-medium">Getting Started</h3>
                <p className="mt-2 text-sm">{placeholderText}</p>
              </div>
              <div className="rounded-md bg-secondary/10 p-4">
                <h3 className="font-medium">Features</h3>
                <p className="mt-2 text-sm">Explore the powerful features available to you.</p>
              </div>
            </div>
          </div>
        </div>
      </SidebarProvider>
    </div>
  );
}
