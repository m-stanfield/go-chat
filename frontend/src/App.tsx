import { useAuth } from "./AuthContext";
import { useEffect, useState, useRef } from "react";
import ServerPage from "./components/ServerPage";
import IconBanner, { IconInfo } from "./components/IconList";
import AuthPage from "./components/AuthPage";
import { fetchUserServers } from "./api/serverApi";

function App() {
  const auth = useAuth();
  const number_of_messages = 20;
  const [showSettings, setShowSettings] = useState(false);
  const settingsRef = useRef<HTMLDivElement>(null);

  // Add click outside handler to close the settings menu
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
    auth.addLogoutCallback(() => {
      console.log("logout callback");
    });
  }, []);

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
        }
      } catch (error) {
        console.error("Error fetching servers:", error);
      }
    };

    fetchServers();
  }, [auth.authState.user]);

  const [server_icons, setServerId] = useState<IconInfo[]>([]);
  const [selectedServerId, setSelectedServerId] = useState<IconInfo>({
    icon_id: -1,
    name: "",
    image_url: undefined,
  });

  function onServerSelect(server_id: number) {
    const selectedServer = server_icons.find((x) => x.icon_id === server_id);
    if (selectedServer !== undefined) {
      setSelectedServerId(selectedServer);
    }
  }
  return (
    <div className="flex-col h-screen w-screen bg-gray-500 flex py-12 px-4 sm:px-6 lg:px-8">
      {auth.authState.isAuthenticated ? (
        <div className="flex flex-col flex-grow overflow-y-auto">
          <div className="flex justify-between items-center mb-4 p-2">
            <IconBanner
              icon_info={server_icons}
              onServerSelect={onServerSelect}
              direction="horizontal"
              selectedIconId={selectedServerId.icon_id}
            />
            <div className="relative" ref={settingsRef}>
              <button
                onClick={() => setShowSettings(!showSettings)}
                className="p-2 text-white hover:bg-gray-600 rounded-full"
              >
                <svg 
                  xmlns="http://www.w3.org/2000/svg" 
                  className="h-6 w-6" 
                  fill="none" 
                  viewBox="0 0 24 24" 
                  stroke="currentColor"
                >
                  <path 
                    strokeLinecap="round" 
                    strokeLinejoin="round" 
                    strokeWidth={2} 
                    d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
                  />
                  <path 
                    strokeLinecap="round" 
                    strokeLinejoin="round" 
                    strokeWidth={2} 
                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                  />
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
            <ServerPage
              server_id={selectedServerId.icon_id}
              number_of_messages={number_of_messages}
            />
          </div>
        </div>
      ) : (
        <AuthPage />
      )}
    </div>
  );
}

export default App;
