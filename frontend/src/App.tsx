import Login from "./components/login";
import { useAuth } from "./AuthContext";
import { useEffect, useState } from "react";
import ServerPage from "./components/ServerPage";
import ServerIconBanner, { ServerInfo } from "./components/ServerIconBanner";

function App() {
  const auth = useAuth();
  useEffect(() => {
    auth.addLogoutCallback(() => {
      console.log("logout callback");
    });
  }, []);
  const servers_ids = [1, 2].map((x) => ({
    server_id: x,
  }));
  const [selectedServerId, setSelectedServerId] = useState<ServerInfo>(
    servers_ids[0],
  );
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
            <ServerPage server_id={selectedServerId.server_id} />
          </div>
        </div>
      ) : (
        <Login />
      )}
    </div>
  );
}

export default App;
