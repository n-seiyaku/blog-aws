import { createContext } from "react";

// 認証コンテキストの型定義
export interface AuthContextType {
  user: { email: string } | null;
  handleLogin: (params: { email: string; password: string }) => Promise<void>;
  handleLogout: () => void;
}

// 認証コンテキストの作成
export const AuthContext = createContext<AuthContextType | undefined>(undefined);
