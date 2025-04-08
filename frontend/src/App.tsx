import { Routes, Route } from "react-router-dom";
import { HomePage } from "./pages/HomePage";
import { ServerPage } from "./pages/ServerPage";
import Login from "./pages/Login";
import SignUp from "./pages/SignUp";
import { ProtectedRoute } from "./components/ProtectedRoute";

function App() {
  return (
    <div className="flex h-screen w-screen bg-background">
      <main className="flex w-full flex-1 flex-col">
        <div className="mx-auto w-full px-2 sm:px-4">
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/login" element={<Login />} />
            <Route path="/signup" element={<SignUp />} />
            <Route 
              path="/servers/:serverId" 
              element={
                <ProtectedRoute>
                  <HomePage />
                </ProtectedRoute>
              } 
            />
          </Routes>
        </div>
      </main>
    </div>
  );
}

export default App;
