// auth.d.ts

export type User = {
    id: string;
    name: string;
    email: string;
};

export type AuthState = {
    isAuthenticated: boolean;
    user: User | null;
};

export type AuthContextType = {
    authState: AuthState;
    login: (user: User) => void;
    logout: () => void;
};
