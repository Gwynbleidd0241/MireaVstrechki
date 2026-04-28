import { createContext, ReactNode, useContext, useState } from "react";

type AuthState = {
  isAuth: boolean;
  userId: number | null;
  email: string | null;
  role: string | null;
  signIn: (token: string, userId: number, email: string, role: string) => void;
  signOut: () => void;
};

const AuthContext = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState(() => localStorage.getItem("token"));
  const [userId, setUserId] = useState<number | null>(() => {
    const v = localStorage.getItem("user_id");
    return v ? Number(v) : null;
  });
  const [email, setEmail] = useState(() => localStorage.getItem("email"));
  const [role, setRole] = useState(() => localStorage.getItem("role"));

  function signIn(t: string, id: number, e: string, r: string) {
    localStorage.setItem("token", t);
    localStorage.setItem("user_id", String(id));
    localStorage.setItem("email", e);
    localStorage.setItem("role", r);
    setToken(t);
    setUserId(id);
    setEmail(e);
    setRole(r);
  }

  function signOut() {
    localStorage.removeItem("token");
    localStorage.removeItem("user_id");
    localStorage.removeItem("email");
    localStorage.removeItem("role");
    setToken(null);
    setUserId(null);
    setEmail(null);
    setRole(null);
  }

  return (
    <AuthContext.Provider
      value={{ isAuth: Boolean(token), userId, email, role, signIn, signOut }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return ctx;
}
