import Login from "./components/login";
import { useAuth } from "./AuthContext";
import { useEffect, useState } from "react";
import ServerPage from "./components/ServerPage";
import ServerIconBanner, { ServerInfo } from "./components/ServerIconBanner";

function App() {
  const auth = useAuth();
  const number_of_messages = 200;
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
          const servers_ids: ServerInfo[] = data.servers.map((server: any) => {
            return {
              server_id: server.ServerId,
              server_name: server.ServerName,
              owner_id: server.OwnerId,
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

  const [servers_ids, setServerId] = useState<ServerInfo[]>([]);
  const [selectedServerId, setSelectedServerId] = useState<ServerInfo>({
    server_id: 0,
    server_name: "",
    owner_id: 0,
  });
  function onServerSelect(server_id: number) {
    const selectedServer = servers_ids.find((x) => x.server_id === server_id);
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
          <ServerIconBanner
            server_ids={servers_ids}
            onServerSelect={onServerSelect}
          />
          <div className="">
            <h1>Server ID: {selectedServerId.server_id}</h1>
          </div>
          <div className="flex flex-col flex-grow overflow-y-auto">
            <ServerPage
              server_id={selectedServerId.server_id}
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
