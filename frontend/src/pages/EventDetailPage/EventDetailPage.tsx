import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { Event, getEvent } from "../../api/events";
import { getTasks, Task } from "../../api/tasks";
import { getParticipants, Participant } from "../../api/participants";
import { AgendaItem, getAgenda } from "../../api/agenda";
import { getUsers, User } from "../../api/users";
import "./EventDetailPage.css";

const statusLabel: Record<Task["status"], string> = {
  todo: "К выполнению",
  in_progress: "В работе",
  done: "Готово",
};

const roleLabel: Record<Participant["role"], string> = {
  participant: "Участник",
  responsible: "Ответственный",
};

export function EventDetailPage() {
  const { id } = useParams<{ id: string }>();
  const eventId = Number(id);

  const [event, setEvent] = useState<Event | null>(null);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [participants, setParticipants] = useState<Participant[]>([]);
  const [agenda, setAgenda] = useState<AgendaItem[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (Number.isNaN(eventId)) {
      setError("Неверный идентификатор встречи");
      setLoading(false);
      return;
    }

    Promise.all([
      getEvent(eventId),
      getTasks(eventId),
      getParticipants(eventId),
      getAgenda(eventId),
      getUsers(),
    ])
      .then(([ev, tk, pp, ag, us]) => {
        setEvent(ev);
        setTasks(tk);
        setParticipants(pp);
        setAgenda(ag);
        setUsers(us);
      })
      .catch((err) =>
        setError(err instanceof Error ? err.message : "Не удалось загрузить встречу"),
      )
      .finally(() => setLoading(false));
  }, [eventId]);

  if (loading) {
    return <div className="event-detail-page">Загрузка...</div>;
  }

  if (error || !event) {
    return (
      <div className="event-detail-page">
        <Link to="/events" className="back-link">
          ← Все мероприятия
        </Link>
        <p className="form-error">{error || "Встреча не найдена"}</p>
      </div>
    );
  }

  const userById = new Map(users.map((u) => [u.id, u]));

  return (
    <div className="event-detail-page">
      <Link to="/events" className="back-link">
        ← Все мероприятия
      </Link>

      <header className="event-detail-header">
        <h1>{event.title}</h1>
        <p className="event-detail-time">
          {new Date(event.start_time).toLocaleString()} —{" "}
          {new Date(event.end_time).toLocaleString()}
        </p>
        <p className="event-detail-description">
          {event.description || "Без описания"}
        </p>
        <p className="event-detail-creator">
          Создал: {userById.get(event.creator_id)?.email ?? `#${event.creator_id}`}
        </p>
      </header>

      <div className="event-detail-grid">
        <section className="panel">
          <h2>Задачи</h2>
          {tasks.length === 0 ? (
            <p className="empty-text">Задач пока нет</p>
          ) : (
            <ul className="task-list">
              {tasks.map((t) => (
                <li key={t.id} className={`task-item task-status-${t.status}`}>
                  <div className="task-item__title">{t.title}</div>
                  <div className="task-item__meta">
                    <span className="task-status-badge">{statusLabel[t.status]}</span>
                    {t.assignee_id != null && (
                      <span>
                        {userById.get(t.assignee_id)?.email ?? `#${t.assignee_id}`}
                      </span>
                    )}
                    {t.due_date && (
                      <span>до {new Date(t.due_date).toLocaleDateString()}</span>
                    )}
                  </div>
                </li>
              ))}
            </ul>
          )}
        </section>

        <section className="panel">
          <h2>Участники</h2>
          {participants.length === 0 ? (
            <p className="empty-text">Пока никого нет</p>
          ) : (
            <ul className="participant-list">
              {participants.map((p) => (
                <li key={p.id} className="participant-item">
                  <span>{userById.get(p.user_id)?.email ?? `#${p.user_id}`}</span>
                  <span className="participant-role">{roleLabel[p.role]}</span>
                </li>
              ))}
            </ul>
          )}
        </section>

        <section className="panel">
          <h2>Повестка</h2>
          {agenda.length === 0 ? (
            <p className="empty-text">Повестка пока пустая</p>
          ) : (
            <ol className="agenda-list">
              {agenda.map((a) => (
                <li
                  key={a.id}
                  className={`agenda-item ${a.is_done ? "agenda-item--done" : ""}`}
                >
                  <div className="agenda-item__title">{a.title}</div>
                  {a.description && (
                    <div className="agenda-item__description">{a.description}</div>
                  )}
                  {a.duration_minutes != null && (
                    <div className="agenda-item__duration">
                      ≈ {a.duration_minutes} мин
                    </div>
                  )}
                </li>
              ))}
            </ol>
          )}
        </section>
      </div>
    </div>
  );
}
