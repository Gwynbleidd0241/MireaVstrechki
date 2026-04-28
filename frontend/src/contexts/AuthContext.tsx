import { createContext, ReactNode, useContext, useState } from "react";

type AuthState = {
  isAuth: boolean;
  email: string | null;
  role: string | null;
  signIn: (token: string, email: string, role: string) => void;
  signOut: () => void;
};

const AuthContext = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState(() => localStorage.getItem("token"));
  const [email, setEmail] = useState(() => localStorage.getItem("email"));
  const [role, setRole] = useState(() => localStorage.getItem("role"));

  function signIn(t: string, e: string, r: string) {
    localStorage.setItem("token", t);
    localStorage.setItem("email", e);
    localStorage.setItem("role", r);
    setToken(t);
    setEmail(e);
    setRole(r);
  }

  function signOut() {
    localStorage.removeItem("token");
    localStorage.removeItem("email");
    localStorage.removeItem("role");
    setToken(null);
    setEmail(null);
    setRole(null);
  }

  return (
    <AuthContext.Provider
      value={{ isAuth: Boolean(token), email, role, signIn, signOut }}
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
