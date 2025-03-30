import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import AuthenticatedLayout from "./components/AuthenticatedLayout";
import ServerPage from "./components/ServerPage";
import Login from "./components/login";
import SignUp from "./components/signup";
import { useAuth } from "./AuthContext";

function App() {
  const auth = useAuth();
  const number_of_messages = 20;

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={
          auth.authState.isAuthenticated ? 
            <Navigate to="/" /> : 
            <Login />
        } />
        <Route path="/signup" element={
          auth.authState.isAuthenticated ? 
            <Navigate to="/" /> : 
            <SignUp />
        } />
        <Route path="/" element={<AuthenticatedLayout />}>
          <Route index element={<Navigate to="/server" />} />
          <Route 
            path="server/:serverId" 
            element={<ServerPage number_of_messages={number_of_messages} />} 
          />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
