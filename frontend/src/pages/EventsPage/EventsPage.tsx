import { useEffect, useState } from "react";
import { createEvent, Event, getEvents } from "../../api/events";
import "./EventsPage.css";

export function EventsPage() {
  const [events, setEvents] = useState<Event[]>([]);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [startTime, setStartTime] = useState("2026-04-27T10:00");
  const [endTime, setEndTime] = useState("2026-04-27T11:00");
  const [error, setError] = useState("");

  async function loadEvents() {
    try {
      const response = await getEvents();
      setEvents(response);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to load events");
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
        start_time: new Date(startTime).toISOString(),
        end_time: new Date(endTime).toISOString(),
      });

      setTitle("");
      setDescription("");
      await loadEvents();
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to create event");
    }
  }

  return (
    <div className="events-page">
      <div className="page-header">
        <div>
          <h1>Рабочие мероприятия</h1>
          <p>Создавайте встречи, фиксируйте задачи и управляйте участниками.</p>
        </div>
      </div>

      <div className="events-grid">
        <section className="panel">
          <h2>Создать мероприятие</h2>

          <form onSubmit={handleCreateEvent} className="event-form">
            <input
              placeholder="Название"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
            />

            <textarea
              placeholder="Описание"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
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

        <section className="panel">
          <h2>Список мероприятий</h2>

          <div className="event-list">
            {events.length === 0 ? (
              <p className="empty-text">Мероприятий пока нет</p>
            ) : (
              events.map((event) => (
                <article className="event-card" key={event.id}>
                  <h3>{event.title}</h3>
                  <p>{event.description || "Без описания"}</p>
                  <span>
                    {new Date(event.start_time).toLocaleString()} —{" "}
                    {new Date(event.end_time).toLocaleString()}
                  </span>
                </article>
              ))
            )}
          </div>
        </section>
      </div>
    </div>
  );
}
