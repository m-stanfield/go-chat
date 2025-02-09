import React, { useState } from "react";

function SignUp() {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [confirmedPassword, setConfirmedPassword] = useState("");
    const [error, setError] = useState("");

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        // Handle login logic here, e.g., send request to backend

        if (password !== confirmedPassword) {
            setError("Passwords do not match");
            return;
        }

        const payload = {
            username: username,
            password: password,
        };

        try {
            // Send POST request to backend
            const response = await fetch("http://localhost:8080/api/user/create", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(payload),
            });

            // Handle response
            if (response.ok) {
                const data = await response.json();
                console.log("Login successful:", data);
            } else {
                setError(response.statusText);
                console.error("Login failed:", response.statusText);
            }
        } catch (error) {
            console.error("Error submitting login:", error);
        }
    };

    return (
        <div className="flex-grow  bg-slate-700 ">
            <form onSubmit={handleSubmit}>
                {error && <div>{error}</div>}
                <div>
                    <label htmlFor="username">Username:</label>
                    <input
                        type="text"
                        id="username"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        required
                    />
                </div>
                <div>
                    <label htmlFor="password">Password:</label>
                    <input
                        type="password"
                        id="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                </div>
                <div>
                    <label htmlFor="confirmpassword">Password:</label>
                    <input
                        type="password"
                        id="confirmpassword"
                        value={confirmedPassword}
                        onChange={(e) => setConfirmedPassword(e.target.value)}
                        required
                    />
                </div>
                <div className="border-black">
                    <button type="submit">Login</button>
                </div>
            </form>
        </div>
    );
}

export default SignUp;
