import { useContext } from "react";
import { AuthContext } from "./AuthContext";

// 認証コンテキストを利用するためのカスタムフック
export function useAuth() {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error("useAuth must be used inside AuthProvider");
  }

  return context;
}
