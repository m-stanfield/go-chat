import React from "react";

export type ServerInfo = {
    server_id: number;
    owner_id: number;
    server_name: string;
};
interface ServerIconBannerProps {
    server_ids: ServerInfo[];
    onServerSelect: (server_id: number) => void;
}
function ServerIconBanner({
    server_ids,
    onServerSelect,
}: ServerIconBannerProps) {
    function onServerSelectGenerator(server_id: number) {
        return (t: React.MouseEvent<HTMLElement>) => {
            t.preventDefault();
            onServerSelect(server_id);
        };
    }
    const items = server_ids.map((s) => (
        <button onClick={onServerSelectGenerator(s.server_id)}>
            Server: {s.server_id}
        </button>
    ));
    return <div className="flex flex-shrink gap-2">{items}</div>;
}

export default ServerIconBanner;
