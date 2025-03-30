import React from "react";

export type IconInfo = {
    icon_id: number;
    name: string;
    image_url: string | undefined;
};

interface IconBannerProps {
    icon_info: IconInfo[];
    onServerSelect: (icon_id: number) => void;
    direction?: "horizontal" | "vertical";  // New prop with default horizontal
}

function IconBanner({ icon_info, onServerSelect, direction = "horizontal" }: IconBannerProps) {
    function onServerSelectGenerator(server_id: number) {
        return (t: React.MouseEvent<HTMLElement>) => {
            t.preventDefault();
            onServerSelect(server_id);
        };
    }

    // Adjust tooltip position based on direction
    const tooltipPosition = direction === "horizontal" 
        ? "top-0 left-20" 
        : "left-full top-0 ml-2";

    const items = icon_info.map((s) => (
        <button
            key={s.icon_id}
            id={s.icon_id.toString()}
            onClick={onServerSelectGenerator(s.icon_id)}
            className="relative group"
        >
            {s.image_url != undefined ? (
                <>
                    <img
                        src={s.image_url}
                        alt={s.name}
                        className="w-16 h-16 rounded-full object-cover object-top"
                    />
                    <div className={`absolute ${tooltipPosition} px-2 py-1 bg-gray-800 text-white text-sm rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap z-10`}>
                        {s.name}
                    </div>
                </>
            ) : (
                <div className="w-16 h-16 rounded-full bg-gray-600 flex items-center justify-center text-white hover:bg-gray-500">
                    {s.name}
                </div>
            )}
        </button>
    ));

    return (
        <div className={`flex ${direction === "horizontal" ? "flex-row" : "flex-col"} gap-2`}>
            {items}
        </div>
    );
}

export default IconBanner;
