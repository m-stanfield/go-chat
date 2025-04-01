import { Toaster } from "sonner";
import { Home, Info } from "lucide-react";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";

export function HomePage() {
  const items = [
    {
      title: "Home",
      url: "/",
      icon: Home,
    },
    {
      title: "Home 2",
      url: "/home",
      icon: Home,
    },
    {
      title: "About",
      url: "/about",
      icon: Info,
    },
  ];
  
  return (
    <div className="flex h-full w-full">
      <Toaster />
      <SidebarProvider>
        <AppSidebar items={items} />
        <div className="flex flex-1 flex-col">
          <div className="sticky top-0 z-10 bg-background px-2 py-4 shadow-sm">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <SidebarTrigger className="ml-1" />
                <h1 className="text-2xl font-bold">My Application</h1>
              </div>
              <div className="flex items-center gap-2">
                {/* Add any right-side elements here if needed */}
              </div>
            </div>
          </div>
          <div className="rounded-lg bg-card p-6 shadow">
            <h1 className="mb-4 text-2xl font-bold text-card-foreground">Home Page</h1>
            <p className="text-muted-foreground">Welcome to the home page of the application.</p>
            <div className="mt-6 grid gap-4 md:grid-cols-2">
              <div className="rounded-md bg-primary/10 p-4">
                <h3 className="font-medium">Getting Started</h3>
                <p className="mt-2 text-sm">
                  Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
                  incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud
                  exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
                  {/* Long lorem ipsum text truncated for brevity */}
                </p>
              </div>
              <div className="rounded-md bg-secondary/10 p-4">
                <h3 className="font-medium">Features</h3>
                <p className="mt-2 text-sm">Explore the powerful features available to you.</p>
              </div>
            </div>
          </div>
        </div>
      </SidebarProvider>
    </div>
  );
} 