import { useState, useEffect, useRef } from "react";
function WebSocketDisplay() {
  const [messages, setMessages] = useState<string[]>([]);
  const [status, setStatus] = useState<string>("Disconnected");

  const ws = useRef<WebSocket | null>(null);
  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8081");

    ws.current = socket;
    socket.onopen = () => {
      setStatus("Connected");
      console.log("WebSocket connection opened.");
    };

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        setMessages((prevMessages) => [
          ...prevMessages,
          JSON.stringify(data, null, 2),
        ]);
      } catch (error) {
        console.error("Error parsing JSON:", error);
      }
    };

    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
      setStatus("Error");
    };

    socket.onclose = () => {
      setStatus("Disconnected");
      console.log("WebSocket connection closed.");
    };

    // Cleanup function to close the socket when the component is unmounted
    return () => {
      socket.close();
    };
  }, []);

  return (
    <div style={{ padding: "20px", fontFamily: "Arial, sans-serif" }}>
      <h1>WebSocket JSON Viewer</h1>
      <p>
        Status: <strong>{status}</strong>
      </p>
      <div
        style={{
          height: "400px",
          overflowY: "scroll",
          border: "1px solid #ccc",
          padding: "10px",
          backgroundColor: "#f9f9f9",
        }}
      >
        {messages.length > 0 ? (
          messages.map((message, index) => (
            <pre
              key={index}
              style={{
                backgroundColor: "#fff",
                padding: "10px",
                margin: "10px 0",
              }}
            >
              {message}
            </pre>
          ))
        ) : (
          <p>No messages received yet.</p>
        )}
      </div>
    </div>
  );
}
// Function that manages input and returns the JSX for rendering
function UseTextInputList(props) {
  const [inputValue, setInputValue] = useState<string>("");
  const [submittedEntries, setSubmittedEntries] = useState<string[]>([]);

  const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(event.target.value);
  };

  const handleSubmit = () => {
    if (inputValue.trim() !== "") {
      setSubmittedEntries([...submittedEntries, inputValue.trim()]);
      setInputValue("");
    }
  };

  return (
    <div>
      <div style={{ display: "flex", gap: "10px", marginBottom: "10px" }}>
        <input
          type="text"
          value={inputValue}
          onChange={handleInputChange}
          placeholder="Enter text here..."
        />
        <button onClick={handleSubmit}>Submit</button>
      </div>
      <h2>Submitted Entries</h2>
      {submittedEntries.length > 0 ? (
        <ul className="text-4xl font-bold text-blue-300 mb-2 ">
          {submittedEntries.map((entry, index) => (
            <li key={index}>{entry}</li>
          ))}
        </ul>
      ) : (
        <p>No entries submitted yet.</p>
      )}
    </div>
  );
}

function CurrentOutput(props) {
  const [count, setCount] = useState(0);
  const [userId, setUserId] = useState(1);
  const [message, setMessage] = useState<string>("");

  const fetchData = () => {
    fetch(`http://localhost:8080/api/user/${userId}`)
      .then((response) => response.text())
      .then((data) => setMessage(data))
      .catch((error) => console.error("Error fetching data:", error));
  };
  return (
    <div className="max-w-md mx-auto space-y-8">
      <div className="text-center">
        <h1 className="text-4xl font-bold text-gray-900 mb-2">
          Welcome to Vite + React
        </h1>
        <p className="text-gray-600">
          Get started by editing{" "}
          <code className="text-sm bg-gray-100 p-1 rounded">src/App.tsx</code>
        </p>
      </div>

      <div className="bg-white p-6 rounded-lg shadow-md">
        <div className="text-center space-y-4">
          <button
            onClick={() => setCount((count) => count + 1)}
            className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-md transition-colors"
          >
            Count is {count}
          </button>

          <button
            onClick={() => setUserId((userId) => userId + 1)}
            className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-md transition-colors"
          >
            Increase User ID
          </button>

          <button
            onClick={fetchData}
            className="block w-full bg-green-500 hover:bg-green-600 text-white font-semibold py-2 px-4 rounded-md transition-colors"
          >
            Fetch from Server
          </button>

          {message && (
            <div className="mt-4 p-4 bg-gray-50 rounded-md">
              <p className="text-gray-700">Server Response:</p>
              <p className="text-gray-900 font-medium">{message}</p>
            </div>
          )}
        </div>
      </div>

      <div className="text-center text-gray-500 text-sm">
        Built with Vite, React, and Tailwind CSS
      </div>
    </div>
  );
}
function App() {
  const [count, setCount] = useState<boolean>(true);
  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div>
        <button
          onClick={() => setCount((count) => !count)}
          className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-md transition-colors"
        >
          {count ? "Hide UI" : "Show UI"}
        </button>
      </div>
      <div>
        <WebSocketDisplay />
      </div>
      <div>
        {count ? (
          <div>
            <CurrentOutput />
          </div>
        ) : (
          <div></div>
        )}
      </div>
    </div>
  );
}

export default App;
