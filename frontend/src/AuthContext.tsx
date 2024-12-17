// context/AuthContext.tsx

import React, { createContext, useContext, useState, ReactNode } from "react";
import { AuthContextType, AuthState, LogoutCallback, User } from "./auth";

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [authState, setAuthState] = useState<AuthState>({
    isAuthenticated: false,
    user: null,
    logoutCallbacks: [],
  });

  const login = (user: User) => {
    setAuthState((prevState) => {
      return {
        isAuthenticated: true,
        user: { ...user },
        logoutCallbacks: [],
      };
    });
  };

  const logout = () => {
    setAuthState((prevState) => {
      prevState.logoutCallbacks.forEach((callback) => callback());
      return {
        isAuthenticated: false,
        user: null,
        logoutCallbacks: [],
      };
    });
  };

  const addLogoutCallback = (callback: LogoutCallback) => {
    const removeCallback = () => {
      setAuthState((prevState) => {
        const cbs_removed = prevState.logoutCallbacks.filter(
          (cb) => cb !== callback,
        );
        return {
          ...prevState,
          logoutCallbacks: cbs_removed,
        };
      });
    };
    setAuthState((prevState) => {
      return {
        ...prevState,
        logoutCallbacks: [callback, ...prevState.logoutCallbacks],
      };
    });

    return removeCallback;
  };
  return (
    <AuthContext.Provider
      value={{ authState, login, logout, addLogoutCallback }}
    >
      {children}
    </AuthContext.Provider>
  );
};

// Custom hook for accessing the auth context
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
