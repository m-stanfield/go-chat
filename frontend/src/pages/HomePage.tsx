import { Toaster } from "sonner";
import { Home } from "lucide-react";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { useEffect, useState } from "react";
import { useAuth } from "@/AuthContext";
import { serverApi, ServerData } from "@/api";
import { useNavigate, useParams } from "react-router-dom";
import ServerPage from "./ServerPage";

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
  const [servers, setServers] = useState<ServerData[]>([]);
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
        const retrievedServers = await serverApi.fetchUserServers(auth.authState.user.id);
        // Set the selected server if serverId is provided
        setServers(retrievedServers);
      } catch (error) {
        console.error("Error fetching servers:", error);
      }
    };

    fetchServers();
  }, [auth.authState.user]);

  useEffect(() => {
    if (!serverId) {
      navigate(`/servers/${servers[0]?.ServerId}`, { replace: true });
    }
    const server = servers.find((s) => s.ServerId === parseInt(serverId || ""));

    if (!server) {
      return;
    }
    setCurrentName(server?.ServerName || "");
    const serverNavItems = servers.map((server: ServerData) => ({
      title: server.ServerName,
      url: `/servers/${server.ServerId}`,
      icon: Home,
      selected: server.ServerId === parseInt(serverId || ""),
    }));
    setServerNavItem(serverNavItems);
    navigate(`/servers/${server.ServerId}`);
  }, [serverId, servers, navigate]);

  return (
    <div className="flex h-full w-full">
      <Toaster />
      <SidebarProvider>
        <AppSidebar items={serverNaveItem} />
        <div className="flex flex-1 flex-col">
          <div className="sticky top-0 z-10 bg-background px-2 py-4 shadow-sm">
            <div className="flex items-center justify-between">
              <SidebarTrigger className="ml-1" />
              <h1 className="text-2xl font-bold">{currentName || ""}</h1>
            </div>
          </div>
            <ServerPage server_id={parseInt(serverId || "")} number_of_messages={25} />
        </div>
      </SidebarProvider>
    </div>
  );
}
