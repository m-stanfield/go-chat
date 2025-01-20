import Login from "./components/login";
import { useAuth } from "./AuthContext";
import { useEffect, useState } from "react";
import ServerPage from "./components/ServerPage";
import IconBanner, { IconInfo } from "./components/IconList";

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
          `http://localhost:8080/api/user/${auth.authState.user?.id}/servers`,
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
          const servers_ids: IconInfo[] = data.servers.map((server: any) => {
            return {
              icon_id: server.ServerId,
              name: server.ServerName,
              image_url:
                "https://miro.medium.com/v2/resize:fit:720/format:webp/0*UD_CsUBIvEDoVwzc.png",
            };
          });
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
    <div className="flex-col h-screen w-screen bg-gray-500 flex py-12 px-4 sm:px-6 lg:px-8 ">
      <button onClick={auth.logout} className="flex w-full">
        Logout
      </button>
      {auth.authState.isAuthenticated ? (
        <div className="flex flex-col flex-grow overflow-y-auto">
          <IconBanner
            icon_info={server_icons}
            onServerSelect={onServerSelect}
          />
          <div className="">
            <h1>Server ID: {selectedServerId.icon_id}</h1>
          </div>
          <div className="flex flex-col flex-grow overflow-y-auto">
            <ServerPage
              server_id={selectedServerId.icon_id}
              number_of_messages={number_of_messages}
            />
          </div>
        </div>
      ) : (
        <Login />
      )}
    </div>
  );
}

export default App;
