import React, { useState } from "react";
import { useAuth } from "../AuthContext.tsx";
import SignUp from "./signup.tsx";

function Login() {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const auth = useAuth();

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        // Handle login logic here, e.g., send request to backend
        const payload = {
            username: username,
            password: password,
        };

        try {
            // Send POST request to backend
            const response = await fetch("http://localhost:8080/api/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include",
                body: JSON.stringify(payload),
            });

            // Handle response
            if (response.ok) {
                const data = await response.json();
                console.log(data);
                auth.login({
                    name: username,
                    id: data.userid,
                    token: data.token,
                    token_expire_time: data.token_expire_time,
                });
                console.log("Login successful:", data);
            } else {
                console.error("Login failed:", response.statusText);
            }
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
