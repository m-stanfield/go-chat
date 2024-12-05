export type MessageData = {
    id: number | undefined;
    author: string;
    author_id: string;
    date: Date;
    message: string;
};

interface MessageProps {
    message: MessageData;
}

function formatDateTime(date: Date) {
    if (date === undefined) {
        return "Unknown";
    }
    return date.getHours() + ":" + date.getMinutes() + ", " + date.toDateString();
}
function Message({ message }: MessageProps) {
    const dayTime = formatDateTime(message.date);
    return (
        <>
            {message.id === undefined ? (
                <div>Invalid Message</div>
            ) : (
                <div className="grid grid-rows-2 gap-1 px-2 py-1">
                    <div className="grid grid-cols-2 gap-1">
                        <div className="col-span-1 overflow-auto text-lg font-bold">
                            {message.author}
                        </div>
                        <div className="col-span-1 overflow-auto text-right text-sm font-thin">
                            {dayTime}
                        </div>
                    </div>
                    <div className="overflow-auto">{message.message}</div>
                </div>
            )}
        </>
    );
}

export default Message;
