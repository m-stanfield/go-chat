import { createContext, useContext, useState, ReactNode, useEffect } from "react";
import { AuthContextType, AuthState, LogoutCallback, User } from "./auth";
import { logoutUser, reconnectSession } from "./api/auth";

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [authState, setAuthState] = useState<AuthState>({
    isAuthenticated: false,
    user: null,
    logoutCallbacks: [],
  });
  const [loading, setLoading] = useState(true);

  // Check for existing session on component mount
  useEffect(() => {
    const checkSession = async () => {
      try {
        const data = await reconnectSession();
        if (data === null) {
          setLoading(false);
          return;
        }
        login({
          id: data.userid,
          name: data.username,
        });
      } catch (error) {
        console.error("Session check failed:", error);
      } finally {
        setLoading(false);
      }
    };

    checkSession();
  }, []);

  const login = (user: User) => {
    setAuthState((prevState) => {
      return {
        isAuthenticated: true,
        user: { ...user },
        logoutCallbacks: [],
      };
    });
  };

  const logout = async () => {
    if (!authState.isAuthenticated) {
      return;
    }

    await logoutUser();

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
        const cbs_removed = prevState.logoutCallbacks.filter((cb) => cb !== callback);
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

  if (loading) {
    return (
      <div className="flex h-full w-full items-center justify-center">
        <div className="h-32 w-32 animate-spin rounded-full border-b-2 border-t-2 border-gray-900 dark:border-white"></div>
      </div>
    );
  }

  return (
    <AuthContext.Provider value={{ authState, login, logout, addLogoutCallback }}>
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
