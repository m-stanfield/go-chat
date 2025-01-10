import Login from "./components/login";
import { useAuth } from "./AuthContext";
import { useEffect } from "react";
import ServerPage from "./components/ServerPage";

function App() {
  const auth = useAuth();
  useEffect(() => {
    auth.addLogoutCallback(() => {
      console.log("logout callback");
    });
  }, []);
  const server_id = 2;
  return (
    <div className="flex-col h-screen w-screen bg-gray-500 py-12 px-4 sm:px-6 lg:px-8 flex">
      <button onClick={auth.logout} className="flex w-full">
        Logout
      </button>
      {auth.authState.isAuthenticated ? (
        <div className=" w-full h-full">
          <div className=" w-full">
            <h1>Server ID: {server_id}</h1>
          </div>
          <ServerPage server_id={server_id} />
        </div>
      ) : (
        <Login />
      )}
    </div>
  );
}

export default App;
