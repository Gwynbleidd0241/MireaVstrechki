import { useState } from "react";
import { register } from "../api/auth";

type Props = {
  onGoLogin: () => void;
};

export function RegisterPage({ onGoLogin }: Props) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState("employee");
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    try {
      await register(email, password, role);
      onGoLogin();
    } catch (err) {
      setError(err instanceof Error ? err.message : "registration failed");
    }
  }

  return (
    <div>
      <h1>Регистрация</h1>

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

        <div>
          <select value={role} onChange={(e) => setRole(e.target.value)}>
            <option value="employee">Сотрудник</option>
            <option value="organizer">Организатор</option>
            <option value="admin">Администратор</option>
          </select>
        </div>

        <button type="submit">Зарегистрироваться</button>
      </form>

      {error && <p style={{ color: "red" }}>{error}</p>}

      <button onClick={onGoLogin}>Уже есть аккаунт</button>
    </div>
  );
}
