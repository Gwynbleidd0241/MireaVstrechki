import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Event, EventStatus, getEvents } from "../../api/events";
import { friendlyError } from "../../api/errors";
import "./CalendarPage.css";

const MONTH_NAMES = [
  "Январь",
  "Февраль",
  "Март",
  "Апрель",
  "Май",
  "Июнь",
  "Июль",
  "Август",
  "Сентябрь",
  "Октябрь",
  "Ноябрь",
  "Декабрь",
];

const DAY_NAMES = ["Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"];

const statusClass: Record<EventStatus, string> = {
  scheduled: "cal-event--scheduled",
  cancelled: "cal-event--cancelled",
  completed: "cal-event--completed",
};

export function CalendarPage() {
  const [events, setEvents] = useState<Event[]>([]);
  const [error, setError] = useState("");
  const [current, setCurrent] = useState(() => {
    const d = new Date();
    return { year: d.getFullYear(), month: d.getMonth() };
  });

  useEffect(() => {
    getEvents()
      .then(setEvents)
      .catch((err) =>
        setError(friendlyError(err, "Не удалось загрузить мероприятия")),
      );
  }, []);

  const today = new Date();
  const { year, month } = current;

  const firstWeekday = (new Date(year, month, 1).getDay() + 6) % 7;
  const daysInMonth = new Date(year, month + 1, 0).getDate();

  const cells: (number | null)[] = [
    ...Array(firstWeekday).fill(null),
    ...Array.from({ length: daysInMonth }, (_, i) => i + 1),
  ];

  while (cells.length % 7 !== 0) cells.push(null);

  function eventsForDay(day: number): Event[] {
    return events.filter((e) => {
      const d = new Date(e.start_time);
      return (
        d.getFullYear() === year &&
        d.getMonth() === month &&
        d.getDate() === day
      );
    });
  }

  function isToday(day: number) {
    return (
      today.getFullYear() === year &&
      today.getMonth() === month &&
      today.getDate() === day
    );
  }

  function prevMonth() {
    setCurrent(({ year, month }) =>
      month === 0 ? { year: year - 1, month: 11 } : { year, month: month - 1 },
    );
  }

  function nextMonth() {
    setCurrent(({ year, month }) =>
      month === 11 ? { year: year + 1, month: 0 } : { year, month: month + 1 },
    );
  }

  function goToday() {
    const d = new Date();
    setCurrent({ year: d.getFullYear(), month: d.getMonth() });
  }

  const monthEvents = events
    .filter((e) => {
      const d = new Date(e.start_time);
      return d.getFullYear() === year && d.getMonth() === month;
    })
    .sort((a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime());

  return (
    <div className="calendar-page">
      <div className="calendar-header">
        <div className="calendar-header__left">
          <h1>
            {MONTH_NAMES[month]} {year}
          </h1>
        </div>
        <div className="calendar-header__controls">
          <button onClick={goToday} className="cal-btn cal-btn--today">
            Сегодня
          </button>
          <button onClick={prevMonth} className="cal-btn cal-btn--nav" aria-label="Предыдущий месяц">‹</button>
          <button onClick={nextMonth} className="cal-btn cal-btn--nav" aria-label="Следующий месяц">›</button>
        </div>
      </div>

      {error && <p className="form-error">{error}</p>}

      <div className="calendar-grid">
        {DAY_NAMES.map((name) => (
          <div key={name} className="calendar-grid__dayname">{name}</div>
        ))}
        {cells.map((day, idx) => (
          <div
            key={idx}
            className={`calendar-cell ${day === null ? "calendar-cell--empty" : ""} ${day && isToday(day) ? "calendar-cell--today" : ""}`}
          >
            {day !== null && (
              <>
                <span className="calendar-cell__num">{day}</span>
                <div className="calendar-cell__events">
                  {eventsForDay(day).map((event) => (
                    <Link
                      key={event.id}
                      to={`/events/${event.id}`}
                      className={`cal-event ${statusClass[event.status]}`}
                      title={`${event.title} — ${new Date(event.start_time).toLocaleTimeString("ru", { hour: "2-digit", minute: "2-digit" })}`}
                    >
                      <span className="cal-event__time">
                        {new Date(event.start_time).toLocaleTimeString("ru", { hour: "2-digit", minute: "2-digit" })}
                      </span>
                      <span className="cal-event__title">{event.title}</span>
                    </Link>
                  ))}
                </div>
              </>
            )}
          </div>
        ))}
      </div>

      <div className="cal-list">
        {monthEvents.length === 0 ? (
          <p className="empty-text">Мероприятий в этом месяце нет</p>
        ) : (
          monthEvents.map((event) => {
            const d = new Date(event.start_time);
            return (
              <Link key={event.id} to={`/events/${event.id}`} className={`cal-list__item cal-list__item--${event.status}`}>
                <div className="cal-list__date">
                  <span className="cal-list__day">{d.getDate()}</span>
                  <span className="cal-list__weekday">{d.toLocaleDateString("ru", { weekday: "short" })}</span>
                </div>
                <div className="cal-list__body">
                  <div className="cal-list__title">{event.title}</div>
                  <div className="cal-list__time">
                    {d.toLocaleTimeString("ru", { hour: "2-digit", minute: "2-digit" })} —{" "}
                    {new Date(event.end_time).toLocaleTimeString("ru", { hour: "2-digit", minute: "2-digit" })}
                  </div>
                </div>
              </Link>
            );
          })
        )}
      </div>
    </div>
  );
}
