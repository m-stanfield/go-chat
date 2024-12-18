import ChatPage from "./ChatPage";

interface ServerPageProps {
    server_id: number;
}
function ServerPage({ server_id }: ServerPageProps) {
    return (
        <div className="w-full h-full">
            <div>
                <h1 className="w-full">Server ID: {server_id}</h1>
            </div>
            <div className="h-full w-full">
                <ChatPage channel_id={2} />
            </div>
        </div>
    );
}

export default ServerPage;
