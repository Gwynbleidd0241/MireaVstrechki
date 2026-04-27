import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { Layout } from "./components/Layout/Layout";
import { EventsPage } from "./pages/EventsPage/EventsPage";
import { LoginPage } from "./pages/LoginPage/LoginPage";
import { RegisterPage } from "./pages/RegisterPage/RegisterPage";

function App() {
  const isAuth = Boolean(localStorage.getItem("token"));

  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/login"
          element={isAuth ? <Navigate to="/events" /> : <LoginPage />}
        />

        <Route
          path="/register"
          element={isAuth ? <Navigate to="/events" /> : <RegisterPage />}
        />

        <Route
          path="/events"
          element={
            isAuth ? (
              <Layout>
                <EventsPage />
              </Layout>
            ) : (
              <Navigate to="/login" />
            )
          }
        />

        <Route
          path="*"
          element={<Navigate to={isAuth ? "/events" : "/login"} />}
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
