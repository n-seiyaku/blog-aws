import { useState } from "react";
import { useAuth } from "../../auth";
import { useNavigate } from "react-router";
import { Button, TextField } from "@mui/material";

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
    <div className="flex flex-col gap-4 w-full max-w-sm px-6">
      <div className="text-center text-lg font-semibold">My Blog</div>

      <TextField
        id="username"
        label="Username"
        variant="outlined"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />

      <TextField
        id="password"
        label="Password"
        variant="outlined"
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />

      <Button variant="contained" onClick={handleSubmit}>
        Login
      </Button>
    </div>
  );
}

export default Login;
