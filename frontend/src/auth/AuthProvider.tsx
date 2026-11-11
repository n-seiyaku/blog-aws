import {useState, type ReactNode} from "react";
import {AuthContext} from "./AuthContext";

interface AuthProviderProps {
  children: ReactNode;
}

type User = {
  id: string;
  email: string;
  role: string;
};

// 認証プロバイダーコンポーネント
export function AuthProvider({children}: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);

  async function handleLogin({email, password}: {email: string; password: string}) {
    const loginUrl: string = import.meta.env.VITE_API_URL + "/api/auth/login";
    const init: RequestInit = {
      method: "POST",
      body: JSON.stringify({
        email: email,
        password: password,
      }),
    };
    // test2@gmail.com
    // 12346578
    try {
      const response = await fetch(loginUrl, init);
      if (response.ok) {
        const {accessToken, refreshToken, user} = await response.json();
        localStorage.setItem("accessToken", accessToken);
        localStorage.setItem("refreshToken", refreshToken);
        console.log("accessToken", accessToken);
        console.log("refreshToken", refreshToken);
        setUser(user);
      }
    } catch (error) {
      console.log(error);
    }
  }

  async function handleLogout() {
    const logoutUrl: string = import.meta.env.VITE_API_URL + "/api/auth/logout";
    const init: RequestInit = {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + localStorage.getItem("accessToken"),
      },
      body: JSON.stringify({
        refreshToken: localStorage.getItem("refreshToken"),
      }),
    };
    try {
      const response = await fetch(logoutUrl, init);
      console.log("response", response);
      if (response.ok) {
        setUser(null);
        localStorage.removeItem("accessToken");
        localStorage.removeItem("refreshToken");
      }
    } catch (error) {
      console.log(error);
    }
  }

  return (
    <AuthContext.Provider value={{user, handleLogin, handleLogout}}>
      {children}
    </AuthContext.Provider>
  );
}
