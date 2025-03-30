import { useAuth } from "./AuthContext";
import { useEffect, useState } from "react";
import ServerPage from "./components/ServerPage";
import IconBanner, { IconInfo } from "./components/IconList";
import AuthPage from "./components/AuthPage";

type ServerIconResponse = {
  ServerId: number;
  ServerName: string;
  image_url: string | undefined;
};
function App() {
  const auth = useAuth();
  const number_of_messages = 20;
  useEffect(() => {
    auth.addLogoutCallback(() => {
      console.log("logout callback");
    });
  }, []);

  useEffect(() => {
    const _call = async () => {
      if (auth.authState.user?.id === undefined) {
        return;
      }

      try {
        // Send POST request to backend
        const response = await fetch(
          `http://localhost:8080/api/users/${auth.authState.user?.id}/servers`,
          {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
            },
            credentials: "include",
          },
        );
        if (response.ok) {
          const data = await response.json();
          const servers_ids: IconInfo[] = data.servers.map(
            (server: ServerIconResponse) => {
              return {
                icon_id: server.ServerId,
                name: server.ServerName,
                image_url:
                  "https://miro.medium.com/v2/resize:fit:720/format:webp/0*UD_CsUBIvEDoVwzc.png",
              };
            },
          );
          setServerId(servers_ids);
          if (servers_ids.length > 0) {
            setSelectedServerId(servers_ids[0]);
          }
        } else {
          console.error("Login failed:", response.statusText);
        }
      } catch (error) {
        console.error("Error submitting login:", error);
        return;
      }
    };
    _call();
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
            <button
              onClick={auth.logout}
              className="flex w-32 justify-center rounded-md bg-red-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-red-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600"
            >
              Logout
            </button>
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
