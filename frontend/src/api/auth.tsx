interface LoginPayload {
    username: string;
    password: string;
}

interface LoginResponse {
    userid: number;
    // Add other response fields if needed
}

interface SessionResponse {
    userid: number;
    username: string;
}

export const loginUser = async (payload: LoginPayload): Promise<LoginResponse> => {
    const response = await fetch("/api/auth/login", {
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

export const logoutUser = async (): Promise<boolean> => {
    // Call the logout API endpoint
    fetch("/api/auth/logout", {
        method: "POST",
        credentials: "include",
    }).catch((error) => {
        console.error("Logout request failed:", error);
        return false;
    });
    return true;
};

export const reconnectSession = async (): Promise<SessionResponse | null> => {
    const response = await fetch("/api/auth/session", {
        method: "POST",
        credentials: "include",
    });

    if (!response.ok) {
        return null;
    }
    const data = await response.json();
    // If session is valid, log the user in
    return {
        userid: data.userid,
        username: data.username,
    };
};
