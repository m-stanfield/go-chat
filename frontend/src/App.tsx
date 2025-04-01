import { Routes, Route } from "react-router-dom";
import { SidebarProvider, SidebarTrigger } from "./components/ui/sidebar";
import { AppSidebar } from "./components/app-sidebar";
function App() {
  return (
    <div className="flex h-screen w-screen bg-gray-100">
      <SidebarProvider>
        <AppSidebar />
        <main>
          <SidebarTrigger />
          <div>Hello</div>
        </main>
      </SidebarProvider>
      <div className="flex-1 p-6 md:ml-64">
        <div className="mx-auto max-w-7xl">
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
    <div className="rounded-lg bg-white p-6 shadow">
      <h1 className="mb-4 text-2xl font-bold">Home Page</h1>
      <p>Welcome to the home page of the application.</p>
    </div>
  );
}

function About() {
  return (
    <div className="rounded-lg bg-white p-6 shadow">
      <h1 className="mb-4 text-2xl font-bold">About Page</h1>
      <p>This is the about page of the application.</p>
    </div>
  );
}

export default App;
