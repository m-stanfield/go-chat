import Login from "./components/login";
import { useAuth } from "./AuthContext";
import ChatPage from "./components/ChatPage";

function App() {
  const auth = useAuth();
  return (
    <div className="h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8 flex">
      <div>
        <button onClick={auth.logout}>Logout</button>
      </div>
      <div>{auth.authState.isAuthenticated ? <ChatPage /> : <Login />}</div>
    </div>
  );
}

export default App;
