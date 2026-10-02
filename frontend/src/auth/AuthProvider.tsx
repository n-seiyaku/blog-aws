import { useState, type ReactNode } from "react";
import { AuthContext } from "./AuthContext";

interface AuthProviderProps {
  children: ReactNode;
}

// 認証プロバイダーコンポーネント
export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<{ email: string } | null>(null);

  async function handleLogin({
    email,
    password,
  }: {
    email: string;
    password: string;
  }) {
    if (email === "123" && password === "123") {
      setUser({ email });
      return;
    }

    throw new Error("Email hoặc mật khẩu không đúng");
  }

  function handleLogout() {
    setUser(null);
  }

  return (
    <AuthContext.Provider value={{ user, handleLogin, handleLogout }}>
      {children}
    </AuthContext.Provider>
  );
}
