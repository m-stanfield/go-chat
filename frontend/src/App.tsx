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
  return (
    <div className="flex-col h-screen w-screen bg-gray-500 py-12 px-4 sm:px-6 lg:px-8 flex">
      <button onClick={auth.logout} className="flex w-full">
        Logout
      </button>
      <div className="flex w-full h-full">
        {auth.authState.isAuthenticated ? (
          <ServerPage server_id={1} />
        ) : (
          <Login />
        )}
      </div>
    </div>
  );
}

export default App;
