import { SyntheticEvent, useState } from "react";

interface MessageSubmitWindowProps {
    onSubmit: (t: SyntheticEvent, inputValue: string) => string;
}
function MessageSubmitWindow({ onSubmit }: MessageSubmitWindowProps) {
    const [inputValue, setInputValue] = useState("");
    const onInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        setInputValue(e.target.value);
    };
    const onKeyDown = (event: React.KeyboardEvent<HTMLTextAreaElement>) => {
        if (event.key === "Enter" && event.shiftKey) {
            event.preventDefault();
            setInputValue(inputValue + "\n");
        } else if (event.key === "Enter") {
            event.preventDefault();
            const newValue = onSubmit(event, inputValue);
            setInputValue(newValue);
        }
    };
    return (
        <div className="flex h-full w-full">
            <textarea
                name="text"
                value={inputValue}
                onChange={onInputChange}
                onKeyDown={onKeyDown}
                className="flex h-full w-full overflow-auto rounded-lg px-2 text-black"
            ></textarea>
        </div>
    );
}
export default MessageSubmitWindow;
