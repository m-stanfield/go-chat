// auth.d.ts

export type LogoutCallback = () => void;
export type LogoutCallbackCancel = () => void;

export type User = {
    id: string;
    name: string;
};

export type AuthState = {
    isAuthenticated: boolean;
    user: User | null;
    logoutCallbacks: LogoutCallback[];
};

export type AuthContextType = {
    authState: AuthState;
    login: (user: User) => void;
    logout: () => void;
    addLogoutCallback: (callback: LogoutCallback) => LogoutCallbackCancel;
};
