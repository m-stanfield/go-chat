import { SyntheticEvent, useState } from "react";

interface MessageSubmitWindowProps {
    onSubmit: (t: SyntheticEvent, inputValue: string) => string;
    validateMessage: (message: string) => string | undefined;
}
function MessageSubmitWindow({ onSubmit, validateMessage }: MessageSubmitWindowProps) {
    const [inputValue, setInputValue] = useState("");
    const [errorMessage, setErrorMessage] = useState("");
    const onInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        const message = e.target.value;
        const error = validateMessage(message);
        if (error) {
            setErrorMessage(error);
        } else {
            setErrorMessage("");
        }
        setInputValue(message);
    };
    const onKeyDown = (event: React.KeyboardEvent<HTMLTextAreaElement>) => {
        if (event.key === "Enter" && event.shiftKey) {
            event.preventDefault();
            setInputValue(inputValue + "\n");
        } else if (event.key === "Enter") {
            event.preventDefault();
            const newValue = onSubmit(event, inputValue);
            setInputValue(newValue);
            if (newValue.length === 0) {
                setErrorMessage("");
            }
        }
    };
    return (
        <div className="flex h-full w-full flex-col">
            <textarea
                name="text"
                value={inputValue}
                onChange={onInputChange}
                onKeyDown={onKeyDown}
                className={`h-full w-full resize-none overflow-auto rounded-lg border px-2 text-black ${errorMessage
                        ? "border-2 border-red-500 focus:outline-none focus:ring-2 focus:ring-red-500"
                        : ""
                    }`}
            ></textarea>
            <div className={`${errorMessage ? "min-h-[20px]" : "min-h-[20px]"}`}>
                {errorMessage && <div className="text-sm text-red-500">{errorMessage}</div>}
            </div>
        </div>
    );
}
export default MessageSubmitWindow;
