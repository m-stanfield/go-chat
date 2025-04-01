import { NavLink } from "react-router-dom";

export function AboutPage() {
  return (
    <div className="rounded-lg bg-card p-6 shadow">
      <NavLink to="/home">Home</NavLink>
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