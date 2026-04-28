import { ReactNode } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import { useAuth } from "../../contexts/AuthContext";
import "./Layout.css";

type Props = {
  children: ReactNode;
};

const roleLabel: Record<string, string> = {
  admin: "Администратор",
  organizer: "Организатор",
  employee: "Сотрудник",
};

export function Layout({ children }: Props) {
  const navigate = useNavigate();
  const { email, role, signOut } = useAuth();

  function handleLogout() {
    signOut();
    navigate("/login");
  }

  return (
    <div className="app-shell">
      <header className="topbar">
        <div className="topbar__inner">
          <div className="topbar__brand">
            <div className="topbar__logo">M</div>
            <span className="topbar__brand-name">MeetFlow</span>
          </div>

          <nav className="topbar__nav">
            <NavLink to="/events">Мероприятия</NavLink>
          </nav>

          <div className="topbar__user">
            <div className="topbar__user-info">
              <span className="topbar__email">{email}</span>
              <span className="topbar__role">
                {role ? roleLabel[role] ?? role : ""}
              </span>
            </div>
            <button onClick={handleLogout} className="topbar__logout">
              Выйти
            </button>
          </div>
        </div>
      </header>

      <main className="app-main">{children}</main>
    </div>
  );
}
