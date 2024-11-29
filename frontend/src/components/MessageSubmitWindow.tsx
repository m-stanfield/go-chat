import { SyntheticEvent, useState } from "react";

interface MessageSubmitWindowProps {
    onSubmit: (t: SyntheticEvent) => void
    onInputChange: (e: React.ChangeEvent<HTMLTextAreaElement>) => void
    inputValue: string
}
function MessageSubmitWindow(props: MessageSubmitWindowProps) {


    return (
        <div className=" flex w-full h-full  ">
            <form onSubmit={props.onSubmit} className="flex w-full h-full">
                <textarea
                    name="text"
                    value={props.inputValue}
                    onChange={props.onInputChange}
                    onKeyPress={(event) => {
                        if (event.key === "Enter") {
                            event.preventDefault();
                            props.onSubmit(event);
                        }
                    }}
                    className="overflow-auto text-black flex w-full h-full"
                ></textarea>
            </form>
        </div>
    )
}
export default MessageSubmitWindow;
