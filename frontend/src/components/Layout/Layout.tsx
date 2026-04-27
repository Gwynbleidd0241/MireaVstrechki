import { ReactNode } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import "./Layout.css";

type Props = {
  children: ReactNode;
};

export function Layout({ children }: Props) {
  const navigate = useNavigate();

  const email = localStorage.getItem("email");
  const role = localStorage.getItem("role");

  function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("email");
    localStorage.removeItem("role");

    navigate("/login");
  }

  return (
    <div className="app-layout">
      <aside className="sidebar">
        <div className="sidebar__brand">
          <div className="sidebar__logo">M</div>
          <div>
            <h2>MeetFlow</h2>
            <p>рабочие встречи</p>
          </div>
        </div>

        <nav className="sidebar__nav">
          <NavLink to="/events">Мероприятия</NavLink>
        </nav>

        <div className="sidebar__profile">
          <p className="sidebar__email">{email}</p>
          <p className="sidebar__role">{role}</p>
          <button onClick={logout}>Выйти</button>
        </div>
      </aside>

      <main className="layout-content">{children}</main>
    </div>
  );
}
