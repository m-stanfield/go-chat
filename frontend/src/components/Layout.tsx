import { ReactNode } from "react";
import { Toaster } from "sonner";
import { SidebarProvider } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { Home, Info, LogIn } from "lucide-react";

interface LayoutProps {
  children: ReactNode;
  showSidebar?: boolean;
}

export function Layout({ children, showSidebar = true }: LayoutProps) {
  const navItems = [
    { title: "Home", url: "/", icon: Home },
    { title: "Login", url: "/login", icon: LogIn },
    { title: "Sign Up", url: "/signup", icon: LogIn },
    { title: "About", url: "/about", icon: Info },
  ];

  if (!showSidebar) {
    return (
      <div className="flex h-full w-full">
        <Toaster />
        {children}
      </div>
    );
  }

  return (
    <div className="flex h-full w-full">
      <Toaster />
      <SidebarProvider>
        <AppSidebar items={navItems} />
        {children}
      </SidebarProvider>
    </div>
  );
} 