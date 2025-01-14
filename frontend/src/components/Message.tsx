export type MessageData = {
    message_id: number;
    channel_id: number;
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
            {message.message_id === undefined ? (
                <div>Invalid Message</div>
            ) : (
                <div className="flex-grow grid grid-rows-2 gap-1 px-2 py-1">
                    <div className="flex-grow  grid grid-cols-2 gap-1">
                        <div className="flex-grow col-span-1 overflow-auto text-lg font-bold">
                            {message.author}
                        </div>
                        <div className="col-span-1 flex-grow overflow-auto text-right text-sm font-thin">
                            {dayTime}
                        </div>
                    </div>
                    <div className="overflow-auto flex flex-grow">{message.message}</div>
                </div>
            )}
        </>
    );
}

export default Message;
