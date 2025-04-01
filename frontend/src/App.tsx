import { Routes, Route } from "react-router-dom";
import { SidebarProvider, SidebarTrigger } from "./components/ui/sidebar";
import { AppSidebar } from "./components/app-sidebar";
import { Toaster } from "sonner";
import { Home } from "lucide-react";

const items = [
  {
    title: "Home",
    url: "/",
    icon: Home,
  },
];
function App() {
  return (
    <div className="flex h-screen w-screen bg-background">
      <SidebarProvider>
        <AppSidebar items={items} />
        <main className="flex-1 overflow-auto">
          <div className="px-2 py-6 sm:px-4">
            <div className="mb-6 flex items-center justify-between">
              <div className="flex items-center gap-2">
                <SidebarTrigger className="ml-1" />
                <h1 className="text-2xl font-bold">My Application</h1>
              </div>
              <div className="flex items-center gap-2">
                {/* Add any right-side elements here if needed */}
              </div>
            </div>
            <div className="mx-auto max-w-4xl px-2 sm:px-4">
              <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/about" element={<About />} />
                {/* Add more routes as needed */}
              </Routes>
            </div>
          </div>
        </main>
      </SidebarProvider>
      <Toaster />
    </div>
  );
}

// Simple page components
function HomePage() {
  return (
    <div className="rounded-lg bg-card p-6 shadow">
      <h1 className="mb-4 text-2xl font-bold text-card-foreground">Home Page</h1>
      <p className="text-muted-foreground">Welcome to the home page of the application.</p>
      <div className="mt-6 grid gap-4 md:grid-cols-2">
        <div className="rounded-md bg-primary/10 p-4">
          <h3 className="font-medium">Getting Started</h3>
          <p className="mt-2 text-sm">
            {" "}
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
            incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud
            exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure
            dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
            Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt
            mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed
            do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
            quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis
            aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
            pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia
            deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipiscing
            elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad
            minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
            consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore
            eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
            qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur
            adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut
            enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea
            commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum
            dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
            culpa qui officia deserunt Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed
            do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
            quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis
            aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
            pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia
            deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipiscing
            elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad
            minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
            consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore
            eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
            qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur
            adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut
            enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea
            commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum
            dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
            culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet,
            consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna
            aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
            aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate
            velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non
            proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum
            dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore
            et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
            laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit
            in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat
            cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
            incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud
            exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure
            dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
            Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt
            mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed
            do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
            quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis
            aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
            pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia
            deserunt mollit anim id est laborum. mollit anim id est laborum.
          </p>
        </div>
        <div className="rounded-md bg-secondary/10 p-4">
          <h3 className="font-medium">Features</h3>
          <p className="mt-2 text-sm">Explore the powerful features available to you.</p>
        </div>
      </div>
    </div>
  );
}

function About() {
  return (
    <div className="rounded-lg bg-card p-6 shadow">
      <h1 className="mb-4 text-2xl font-bold text-card-foreground">About Page</h1>
      <p className="mb-4 text-muted-foreground">This is the about page of the application.</p>
      <div className="prose max-w-none">
        <p>
          Our application is designed to provide a seamless user experience with a modern interface.
          We've built this using React, TypeScript, and Tailwind CSS with the shadcn/ui component
          library.
        </p>
        <p className="mt-4">
          The sidebar navigation makes it easy to access different sections of the application, and
          the responsive design ensures a great experience on all devices.
        </p>
      </div>
    </div>
  );
}

export default App;
