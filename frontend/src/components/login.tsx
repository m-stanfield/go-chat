import React, { useState } from "react";
import { useAuth } from "../AuthContext.tsx";
import SignUp from "./signup.tsx";
import { loginUser } from "../api/auth";

function Login() {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const auth = useAuth();

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        
        try {
            const data = await loginUser({
                username,
                password,
            });
            
            auth.login({
                name: username,
                id: data.userid,
            });
            console.log("Login successful:", data);
        } catch (error) {
            console.error("Error submitting login:", error);
        }
    };

    return (
        <div className="flex-grow  bg-slate-700 ">
            <form onSubmit={handleSubmit}>
                <div>
                    <label htmlFor="username">Username:</label>
                    <input
                        type="text"
                        id="username"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                    />
                </div>
                <div>
                    <label htmlFor="password">Password:</label>
                    <input
                        type="password"
                        id="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                </div>
                <div className="border-black">
                    <button type="submit">Login</button>
                </div>
            </form>
            <div>
                <SignUp />
            </div>
        </div>
    );
}

export default Login;
