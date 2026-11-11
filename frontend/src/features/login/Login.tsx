import { useState } from "react";
import { useAuth } from "../../auth";
import { Button, TextField } from "@mui/material";

function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const { handleLogin } = useAuth();

  return (
    <div className="flex flex-col gap-4 w-full max-w-sm px-6">
      <div className="text-center text-lg font-semibold">My Blog</div>

      <TextField
        id="username"
        label="Username"
        variant="outlined"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />

      <TextField
        id="password"
        label="Password"
        variant="outlined"
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />

      <Button
        variant="contained"
        onClick={() => handleLogin({ email, password })}
      >
        Login
      </Button>
    </div>
  );
}

export default Login;
