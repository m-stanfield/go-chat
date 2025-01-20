import React from "react";

export type IconInfo = {
    icon_id: number;
    name: string;
    image_url: string | undefined;
};
interface IconBannerProps {
    icon_info: IconInfo[];
    onServerSelect: (icon_id: number) => void;
}
function IconBanner({ icon_info, onServerSelect }: IconBannerProps) {
    function onServerSelectGenerator(server_id: number) {
        return (t: React.MouseEvent<HTMLElement>) => {
            t.preventDefault();
            onServerSelect(server_id);
        };
    }
    const items = icon_info.map((s) => (
        <button
            id={s.icon_id.toString()}
            onClick={onServerSelectGenerator(s.icon_id)}
        >
            {s.image_url != undefined ? (
                <img
                    src={s.image_url}
                    alt={s.name}
                    className="w-16 h-16 rounded-full object-cover object-top"
                />
            ) : (
                <div> {s.name}</div>
            )}
        </button>
    ));
    return <div className="flex flex-shrink gap-2">{items}</div>;
}

export default IconBanner;
