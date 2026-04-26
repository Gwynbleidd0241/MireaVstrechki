type Props = {
  onLogout: () => void;
};

export function EventsPage({ onLogout }: Props) {
  const email = localStorage.getItem("email");
  const role = localStorage.getItem("role");

  return (
    <div>
      <h1>Рабочие мероприятия</h1>

      <p>
        Пользователь: {email} ({role})
      </p>

      <button onClick={onLogout}>Выйти</button>
    </div>
  );
}
