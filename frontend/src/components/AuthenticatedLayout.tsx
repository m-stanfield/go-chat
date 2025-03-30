import { Navigate, Outlet, useNavigate, useParams } from "react-router-dom";
import { useAuth } from "../AuthContext";
import { useEffect, useRef, useState } from "react";
import IconBanner, { IconInfo } from "./IconList";
import { fetchUserServers } from "../api/serverApi";

function AuthenticatedLayout() {
  const auth = useAuth();
  const navigate = useNavigate();
  const { serverId } = useParams();
  const [showSettings, setShowSettings] = useState(false);
  const settingsRef = useRef<HTMLDivElement>(null);
  const [server_icons, setServerId] = useState<IconInfo[]>([]);
  const [selectedServerId, setSelectedServerId] = useState<IconInfo>({
    icon_id: -1,
    name: "",
    image_url: undefined,
  });

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (settingsRef.current && !settingsRef.current.contains(event.target as Node)) {
        setShowSettings(false);
      }
    }

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [settingsRef]);

  useEffect(() => {
    const fetchServers = async () => {
      if (auth.authState.user?.id === undefined) {
        return;
      }

      try {
        const servers_ids = await fetchUserServers(auth.authState.user.id);
        setServerId(servers_ids);
        if (servers_ids.length > 0) {
          setSelectedServerId(servers_ids[0]);
          navigate(`/server/${servers_ids[0].icon_id}`);
        }
      } catch (error) {
        console.error("Error fetching servers:", error);
      }
    };

    fetchServers();
  }, [auth.authState.user, navigate]);

  useEffect(() => {
    if (serverId && server_icons.length > 0) {
      const id = parseInt(serverId);
      const server = server_icons.find(s => s.icon_id === id);
      if (server) {
        setSelectedServerId(server);
      }
    }
  }, [serverId, server_icons]);

  if (!auth.authState.isAuthenticated) {
    return <Navigate to="/login" />;
  }

  return (
    <div className="flex-col h-screen w-screen bg-gray-500 flex py-12 px-4 sm:px-6 lg:px-8">
      <div className="flex flex-col flex-grow overflow-y-auto">
        <div className="flex justify-between items-center mb-4 p-2">
          <IconBanner
            icon_info={server_icons}
            onServerSelect={(id) => {
              navigate(`/server/${id}`);
            }}
            direction="horizontal"
            selectedIconId={selectedServerId.icon_id}
          />
          <div className="relative" ref={settingsRef}>
            <button
              onClick={() => setShowSettings(!showSettings)}
              className="p-2 text-white hover:bg-gray-600 rounded-full"
            >
              {/* Settings SVG */}
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                {/* ... SVG paths ... */}
              </svg>
            </button>
            {showSettings && (
              <div className="absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-gray-700 ring-1 ring-black ring-opacity-5">
                <div className="py-1" role="menu" aria-orientation="vertical">
                  <button
                    onClick={() => {
                      setShowSettings(false);
                      auth.logout();
                    }}
                    className="block w-full text-left px-4 py-2 text-sm text-white hover:bg-gray-600"
                    role="menuitem"
                  >
                    Logout
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
        <div className="flex flex-col flex-grow overflow-y-auto">
          <Outlet />
        </div>
      </div>
    </div>
  );
}

export default AuthenticatedLayout; 