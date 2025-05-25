import React, {
    createContext,
    useContext,
    useEffect,
    useRef,
    useState,
    ReactNode,
} from 'react';

interface WebSocketMessage {
    message_type: string;
    payload: Record<string, unknown>;
}

interface IWebSocketContext {
    sendMessage: (message: WebSocketMessage) => void;
    ready: boolean;
}

const WebSocketContext = createContext<IWebSocketContext | undefined>(undefined);

interface WebSocketProviderProps {
    url: string;
    onMessage?: (event: MessageEvent) => void;
    children: ReactNode;
}

export const WebSocketProvider: React.FC<WebSocketProviderProps> = ({ url, onMessage, children }) => {
    const socketRef = useRef<WebSocket | null>(null);
    const [ready, setReady] = useState<boolean>(false);

    useEffect(() => {
        const socket = new WebSocket(url);
        socketRef.current = socket;

        socket.onopen = () => {
            setReady(true);
            console.log('WebSocket connected');
        };

        socket.onmessage = (event: MessageEvent) => {
            console.log('WebSocket message received:', event.data);
            if (onMessage) {
                onMessage(event);
            }
        };

        socket.onclose = () => {
            setReady(false);
            console.log('WebSocket disconnected');
        };

        socket.onerror = (error: Event) => {
            console.error('WebSocket error', error);
        };

        return () => {
            socket.close();
        };
    }, [url, onMessage]);

    const sendMessage = (message: WebSocketMessage) => {
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            const stringifiedMessage = JSON.stringify(message);
            socketRef.current.send(stringifiedMessage);
        } else {
            console.warn('WebSocket not ready to send');
        }
    };

    return (
        <WebSocketContext.Provider value={{ sendMessage, ready }}>
            {children}
        </WebSocketContext.Provider>
    );
};

export const useWebSocket = (): IWebSocketContext => {
    const context = useContext(WebSocketContext);
    if (!context) {
        throw new Error('useWebSocket must be used within a WebSocketProvider');
    }
    return context;
};

