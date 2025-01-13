import React from "react";

export type ServerInfo = {
    server_id: number;
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
        <div key={s.server_id} className="w-full">
            <button onClick={onServerSelectGenerator(s.server_id)}>
                Server: ${s.server_id}
            </button>
        </div>
    ));
    return <div>{items}</div>;
}

export default ServerIconBanner;
