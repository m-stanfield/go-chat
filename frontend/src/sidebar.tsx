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
            <div className="  relative flex h-full w-full flex-col bg-gray-800 text-white">
                <nav className="flex-col flex  space-y-2 p-2">
                    <button
                        className="mb-6 flex h-12 w-12 items-center justify-center text-2xl focus:outline-none"
                        onClick={toggleSidebar}
                        aria-label={isOpen ? "Close Sidebar" : "Open Sidebar"}
                    >
                        {isOpen ? "✕" : "☰"}
                    </button>
                    {navItems.map((item) => (
                        <a
                            key={item.label}
                            href={item.href}
                            className=" h-12 flex items-center rounded-md p-2 transition-colors duration-200 hover:bg-gray-700"
                        >
                            <div className=" grid-cols-2">
                                <span className=" mr-2  w-96 items-center justify-left text-xl">
                                    {item.icon}
                                </span>
                                <span className={` ${isOpen ? "inline-block" : "hidden"} `}>
                                    {item.label}
                                </span>
                            </div>
                        </a>
                    ))}
                </nav>
            </div>
        </div>
    );
};

export default Sidebar;
