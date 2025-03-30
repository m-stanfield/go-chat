interface LoginPayload {
    username: string;
    password: string;
}

interface LoginResponse {
    userid: number;
    // Add other response fields if needed
}

export const loginUser = async (payload: LoginPayload): Promise<LoginResponse> => {
    const response = await fetch("http://localhost:8080/api/auth/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(payload),
    });

    if (!response.ok) {
        throw new Error(response.statusText);
    }

    return response.json();
}; 