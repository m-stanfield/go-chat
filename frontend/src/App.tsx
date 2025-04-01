import { Routes, Route } from "react-router-dom";
import Sidebar from "./components/Sidebar";

function App() {
  return (
    <div className="flex h-screen w-screen bg-gray-100">
      <Sidebar />
      <div className="flex-1 md:ml-64 p-6">
        <div className="max-w-7xl mx-auto">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/about" element={<About />} />
            {/* Add more routes as needed */}
          </Routes>
        </div>
      </div>
    </div>
  );
}

// Simple page components
function Home() {
  return (
    <div className="bg-white p-6 rounded-lg shadow">
      <h1 className="text-2xl font-bold mb-4">Home Page</h1>
      <p>Welcome to the home page of the application.</p>
    </div>
  );
}

function About() {
  return (
    <div className="bg-white p-6 rounded-lg shadow">
      <h1 className="text-2xl font-bold mb-4">About Page</h1>
      <p>This is the about page of the application.</p>
    </div>
  );
}

export default App;
