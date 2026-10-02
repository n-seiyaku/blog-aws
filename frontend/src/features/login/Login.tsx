import { useState } from "react";
import { useAuth } from "../../auth";
import { useNavigate } from "react-router";

function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const { user, handleLogin } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async () => {
    await handleLogin({ email: username, password });
    if (user) {
      navigate("/blog");
    }

    // setUsername("");
    // setPassword("");
  };
  return (
    <div>
      <label htmlFor="username">Username</label>
      <input
        type="text"
        id="username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <br />
      <label htmlFor="password">Password</label>
      <input
        type="password"
        id="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <br />
      <button onClick={handleSubmit}>Login</button>
    </div>
  );
}

export default Login;
