import { useState } from "react";
import { login } from "../api/auth";

type Props = {
  onLogin: () => void;
  onGoRegister: () => void;
};

export function LoginPage({ onLogin, onGoRegister }: Props) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    try {
      const response = await login(email, password);

      if (!response.token) {
        throw new Error("token not found");
      }

      localStorage.setItem("token", response.token);
      localStorage.setItem("email", response.email);
      localStorage.setItem("role", response.role);

      onLogin();
    } catch (err) {
      setError(err instanceof Error ? err.message : "login failed");
    }
  }

  return (
    <div>
      <h1>Вход</h1>

      <form onSubmit={handleSubmit}>
        <div>
          <input
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </div>

        <div>
          <input
            placeholder="Пароль"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>

        <button type="submit">Войти</button>
      </form>

      {error && <p style={{ color: "red" }}>{error}</p>}

      <button onClick={onGoRegister}>Создать аккаунт</button>
    </div>
  );
}
