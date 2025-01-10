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
                <div className="w-full grid grid-rows-2 gap-1 px-2 py-1">
                    <div className="w-full grid grid-cols-2 gap-1">
                        <div className="w-full col-span-1 overflow-auto text-lg font-bold">
                            {message.author}
                        </div>
                        <div className="col-span-1 w-full overflow-auto text-right text-sm font-thin">
                            {dayTime}
                        </div>
                    </div>
                    <div className="overflow-auto w-full">{message.message}</div>
                </div>
            )}
        </>
    );
}

export default Message;
