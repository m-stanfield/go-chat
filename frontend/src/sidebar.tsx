import { useState } from "react";

type NavItem = {
    icon: string;
    label: string;
    href: string;
};

export interface NavItemProp {
    items: NavItem[];
    min_width: string;
}

const Sidebar: React.FC<NavItemProp> = ({ items, min_width }) => {
    const [isOpen, setIsOpen] = useState(false);

    const toggleSidebar = () => setIsOpen(!isOpen);
    const width = "w-" + min_width;

    return (
        <div
            className={`fixed left-0 top-0 z-40 flex h-screen transition-all duration-300 ease-in-out ${isOpen ? "w-64" : width
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
                    {items.map((item: NavItem) => (
                        <a
                            key={item.label}
                            href={item.href}
                            className="h-12 flex items-center rounded-md p-2 transition-colors duration-200 hover:bg-gray-700 "
                        >
                            <div className="flex w-full items-center">
                                {/* First span: icon */}
                                <span className="w-10 flex-shrink-0 flex items-center justify-center text-xl">
                                    {item.icon}
                                </span>

                                {/* Second span: label */}
                                <span
                                    className={`${isOpen ? "inline-block" : "hidden"} ml-4 text-base`}
                                >
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
export type { NavItem as NavItemType };
export default Sidebar;
