import { SyntheticEvent, useState } from "react";

interface MessageSubmitWindowProps {
    onSubmit: (t: SyntheticEvent, inputValue: string) => string;
}
function MessageSubmitWindow(props: MessageSubmitWindowProps) {
    const [inputValue, setInputValue] = useState("");
    const onInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        setInputValue(e.target.value);
    };
    return (
        <div className=" flex w-full h-full  ">
            <textarea
                name="text"
                value={inputValue}
                onChange={onInputChange}
                onKeyPress={(event) => {
                    if (event.key === "Enter") {
                        event.preventDefault();
                        const newValue = props.onSubmit(event, inputValue);
                        setInputValue(newValue);
                    }
                }}
                className="overflow-auto text-black flex w-full h-full"
            ></textarea>
        </div>
    );
}
export default MessageSubmitWindow;
