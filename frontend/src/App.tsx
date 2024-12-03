import { useState } from "react";
import Login from "./components/login";
import { AuthProvider, useAuth } from "./AuthContext.tsx";

function App() {
  const [count, setCount] = useState<boolean>(true);
  const auth = useAuth();
  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div>
        <button
          onClick={() => setCount((count) => !count)}
          className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-md transition-colors"
        >
          {count ? "Hide UI" : "Show UI"}
        </button>
      </div>
      <div>
        {auth.authState.isAuthenticated ? <div> Authed!</div> : <Login />}
      </div>
    </div>
  );
}

export default App;
