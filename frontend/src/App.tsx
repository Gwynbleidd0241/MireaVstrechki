import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { Layout } from "./components/Layout/Layout";
import { AuthProvider, useAuth } from "./contexts/AuthContext";
import { EventDetailPage } from "./pages/EventDetailPage/EventDetailPage";
import { EventsPage } from "./pages/EventsPage/EventsPage";
import { LoginPage } from "./pages/LoginPage/LoginPage";
import { RegisterPage } from "./pages/RegisterPage/RegisterPage";

function AppRoutes() {
  const { isAuth } = useAuth();

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
          path="/events/:id"
          element={
            isAuth ? (
              <Layout>
                <EventDetailPage />
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

function App() {
  return (
    <AuthProvider>
      <AppRoutes />
    </AuthProvider>
  );
}

export default App;
