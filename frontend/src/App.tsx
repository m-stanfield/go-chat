import { Routes, Route } from "react-router-dom";
import { HomePage } from "./pages/HomePage";
import { AboutPage } from "./pages/AboutPage";

function App() {
  return (
    <div className="flex h-screen w-screen bg-background">
      <main className="flex w-full flex-1 flex-col">
          <div className="mx-auto w-full px-2 sm:px-4">
            <Routes>
              <Route path="/" element={<HomePage />} />
              <Route path="/home" element={<HomePage />} />
              <Route path="/about" element={<AboutPage />} />
              {/* Add more routes as needed */}
            </Routes>
          </div>
      </main>
    </div>
  );
}

export default App;
