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
    setAuthState({
      isAuthenticated: true,
      user,
      logoutCallbacks: [],
    });
  };

  const logout = () => {
    authState.logoutCallbacks.forEach((callback) => callback());

    setAuthState({
      isAuthenticated: false,
      user: null,
      logoutCallbacks: [],
    });
  };

  const addLogoutCallback = (callback: LogoutCallback) => {
    const removeCallback = () => {
      const cbs_removed = authState.logoutCallbacks.filter(
        (cb) => cb !== callback,
      );
      if (cbs_removed.length == 0) {
        return;
      }

      setAuthState({
        isAuthenticated: authState.isAuthenticated,
        user: authState.user,
        logoutCallbacks: cbs_removed,
      });
    };
    setAuthState({
      isAuthenticated: authState.isAuthenticated,
      user: authState.user,
      logoutCallbacks: [callback, ...authState.logoutCallbacks],
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
