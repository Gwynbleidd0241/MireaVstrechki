import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { register } from "../../api/auth";
import { friendlyError } from "../../api/errors";
import "../LoginPage/LoginPage.css";

export function RegisterPage() {
  const navigate = useNavigate();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState("employee");
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    try {
      await register(email, password, role);
      navigate("/login");
    } catch (err) {
      setError(friendlyError(err, "Не удалось зарегистрироваться"));
    }
  }

  return (
    <div className="auth-page">
      <div className="auth-card">
        <div className="auth-card__header">
          <h1>Регистрация</h1>
          <p>Создайте аккаунт для работы с мероприятиями</p>
        </div>

        <form onSubmit={handleSubmit} className="auth-form">
          <input
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />

          <input
            placeholder="Пароль"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />

          <select value={role} onChange={(e) => setRole(e.target.value)}>
            <option value="employee">Сотрудник</option>
            <option value="organizer">Организатор</option>
            <option value="admin">Администратор</option>
          </select>

          <button type="submit">Зарегистрироваться</button>
        </form>

        {error && <p className="auth-error">{error}</p>}

        <p className="auth-link">
          Уже есть аккаунт? <Link to="/login">Войти</Link>
        </p>
      </div>
    </div>
  );
}
