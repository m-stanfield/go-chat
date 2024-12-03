import React, { useState } from "react";

interface NavItem {
    icon: string;
    label: string;
    href: string;
}

const navItems: NavItem[] = [
    { icon: "🏠", label: "Home", href: "#" },
    { icon: "⚙️", label: "Settings", href: "#" },
    { icon: "❓", label: "Help", href: "#" },
];

const Sidebar: React.FC = () => {
    const [isOpen, setIsOpen] = useState(false);

    const toggleSidebar = () => setIsOpen(!isOpen);

    return (
        <div
            className={`fixed left-0 top-0 z-40 flex h-screen transition-all duration-300 ease-in-out ${isOpen ? "w-64" : "w-16"
                }`}
        >
            <div className="flex h-full w-full flex-col bg-gray-800 text-white">
                <button
                    className="absolute right-4 top-4 text-white focus:outline-none"
                    onClick={toggleSidebar}
                    aria-label={isOpen ? "Close Sidebar" : "Open Sidebar"}
                >
                    {isOpen ? "✕" : "☰"}
                </button>
                <nav className="mt-16 flex flex-col space-y-2 p-2">
                    {navItems.map((item) => (
                        <a
                            key={item.label}
                            href={item.href}
                            className="flex items-center rounded-md p-2 transition-colors duration-200 hover:bg-gray-700"
                        >
                            <span className="mr-2 text-xl">{item.icon}</span>
                            <span className={`${isOpen ? "inline-block" : "hidden"}`}>
                                {item.label}
                            </span>
                        </a>
                    ))}
                </nav>
            </div>
        </div>
    );
};

export default Sidebar;
