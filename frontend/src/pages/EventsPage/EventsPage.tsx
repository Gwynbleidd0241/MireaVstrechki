import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { createEvent, Event, EventStatus, getEvents } from "../../api/events";
import { friendlyError } from "../../api/errors";
import { useAuth } from "../../contexts/AuthContext";
import "./EventsPage.css";

const statusLabel: Record<EventStatus, string> = {
  scheduled: "Запланировано",
  cancelled: "Отменено",
  completed: "Завершено",
};

const statusClass: Record<EventStatus, string> = {
  scheduled: "status--scheduled",
  cancelled: "status--cancelled",
  completed: "status--completed",
};

export function EventsPage() {
  const { role } = useAuth();
  const canCreate = role === "admin" || role === "organizer";

  const [events, setEvents] = useState<Event[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [location, setLocation] = useState("");
  const [meetingUrl, setMeetingUrl] = useState("");
  const [startTime, setStartTime] = useState(() => {
    const d = new Date();
    d.setHours(10, 0, 0, 0);
    return d.toISOString().slice(0, 16);
  });
  const [endTime, setEndTime] = useState(() => {
    const d = new Date();
    d.setHours(11, 0, 0, 0);
    return d.toISOString().slice(0, 16);
  });
  const [error, setError] = useState("");

  async function loadEvents() {
    try {
      const response = await getEvents();
      setEvents(response);
    } catch (err) {
      setError(friendlyError(err, "Не удалось загрузить мероприятия"));
    }
  }

  useEffect(() => {
    loadEvents();
  }, []);

  async function handleCreateEvent(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    try {
      await createEvent({
        title,
        description,
        location,
        meeting_url: meetingUrl,
        start_time: new Date(startTime).toISOString(),
        end_time: new Date(endTime).toISOString(),
      });

      setTitle("");
      setDescription("");
      setLocation("");
      setMeetingUrl("");
      setShowForm(false);
      await loadEvents();
    } catch (err) {
      setError(friendlyError(err, "Не удалось создать мероприятие"));
    }
  }

  return (
    <div className="events-page">
      <div className="page-header">
        <h1>Мероприятия</h1>
        {canCreate && (
          <button className="create-btn" onClick={() => setShowForm((v) => !v)}>
            {showForm ? "Отмена" : "+ Создать"}
          </button>
        )}
      </div>

      <div className={`events-grid ${canCreate && showForm ? "" : "events-grid--single"}`}>
        {canCreate && showForm && (
          <section className="panel">
            <h2>Новое мероприятие</h2>

            <form onSubmit={handleCreateEvent} className="event-form">
              <input
                placeholder="Название"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
              />

              <textarea
                placeholder="Описание"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />

              <input
                placeholder="Место проведения (необязательно)"
                value={location}
                onChange={(e) => setLocation(e.target.value)}
              />

              <input
                placeholder="Ссылка на онлайн-встречу (необязательно)"
                value={meetingUrl}
                onChange={(e) => setMeetingUrl(e.target.value)}
              />

              <label>
                Начало
                <input
                  type="datetime-local"
                  value={startTime}
                  onChange={(e) => setStartTime(e.target.value)}
                />
              </label>

              <label>
                Окончание
                <input
                  type="datetime-local"
                  value={endTime}
                  onChange={(e) => setEndTime(e.target.value)}
                />
              </label>

              <button type="submit">Создать встречу</button>
            </form>

            {error && <p className="form-error">{error}</p>}
          </section>
        )}

        <section className="panel">
          <h2>Список мероприятий</h2>

          <div className="event-list">
            {events.length === 0 ? (
              <p className="empty-text">Мероприятий пока нет</p>
            ) : (
              events.map((event) => (
                <Link
                  to={`/events/${event.id}`}
                  className="event-card-link"
                  key={event.id}
                >
                  <article className="event-card">
                    <div className="event-card__header">
                      <h3>{event.title}</h3>
                      <span
                        className={`status-badge ${statusClass[event.status]}`}
                      >
                        {statusLabel[event.status]}
                      </span>
                    </div>
                    <p>{event.description || "Без описания"}</p>
                    <span>
                      {new Date(event.start_time).toLocaleString()} —{" "}
                      {new Date(event.end_time).toLocaleString()}
                    </span>
                    {event.location && (
                      <span className="event-card__location">
                        📍 {event.location}
                      </span>
                    )}
                  </article>
                </Link>
              ))
            )}
          </div>

          {!canCreate && error && <p className="form-error">{error}</p>}
        </section>
      </div>
    </div>
  );
}
