import { Routes, Route } from "react-router-dom";
import { HomePage } from "./pages/HomePage";
import Login from "./pages/Login";
import SignUp from "./pages/SignUp";
import { ProtectedRoute } from "./components/ProtectedRoute";

function App() {
  return (
    <main className="h-screen w-screen bg-background">
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<SignUp />} />
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <HomePage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/servers/:serverId"
          element={
            <ProtectedRoute>
              <HomePage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/servers/:serverId/channels/:channelId"
          element={
            <ProtectedRoute>
              <HomePage />
            </ProtectedRoute>
          }
        />
      </Routes>
    </main>
  );
}

export default App;
