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
                <div className="">
                    <div className="">
                        <div className="">{message.author}</div>
                        <div className="">{dayTime}</div>
                    </div>
                    <div className="">{message.message}</div>
                </div>
            )}
        </>
    );
}

export default Message;
