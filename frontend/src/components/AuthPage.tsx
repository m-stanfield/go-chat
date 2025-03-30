import { useState } from "react";
import Login from "./login";
import SignUp from "./signup";

function AuthPage() {
    const [showSignUp, setShowSignUp] = useState(false);

    return showSignUp ? (
        <SignUp onSwitchToLogin={() => setShowSignUp(false)} />
    ) : (
        <Login onSwitchToSignUp={() => setShowSignUp(true)} />
    );
}

export default AuthPage; 