import { useState } from "react";
import { LoginPage } from "./pages/LoginPage";
import { RegisterPage } from "./pages/RegisterPage";
import { EventsPage } from "./pages/EventsPage";

type Page = "login" | "register" | "events";

function App() {
  const [page, setPage] = useState<Page>(
    localStorage.getItem("token") ? "events" : "login",
  );

  function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("email");
    localStorage.removeItem("role");
    setPage("login");
  }

  if (page === "register") {
    return <RegisterPage onGoLogin={() => setPage("login")} />;
  }

  if (page === "events") {
    return <EventsPage onLogout={logout} />;
  }

  return (
    <LoginPage
      onLogin={() => setPage("events")}
      onGoRegister={() => setPage("register")}
    />
  );
}

export default App;
