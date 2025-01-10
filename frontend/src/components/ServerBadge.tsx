import React from "react";

interface IconButtonProps {
    server_id: number;
    callback: (server_id: number) => void;
}
function IconButton({ server_id, callback }: IconButtonProps) {
    function onClick(t: React.MouseEvent<HTMLElement>) {
        t.preventDefault();
        callback(server_id);
    }
    return (
        <div>
            <button onClick={onClick}>Server ID: ${server_id}</button>
        </div>
    );
}

export default IconButton;
