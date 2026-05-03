import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { Event, EventStatus, getEvent } from "../../api/events";
import { friendlyError } from "../../api/errors";
import {
  createTask,
  deleteTask,
  getTasks,
  Task,
  updateTask,
} from "../../api/tasks";
import {
  addParticipant,
  getParticipants,
  Participant,
  removeParticipant,
  RSVPStatus,
  updateParticipant,
  updateRSVP,
} from "../../api/participants";
import {
  AgendaItem,
  createAgendaItem,
  deleteAgendaItem,
  getAgenda,
  updateAgendaItem,
} from "../../api/agenda";
import { getUsers, User } from "../../api/users";
import { useAuth } from "../../contexts/AuthContext";
import "./EventDetailPage.css";

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

const taskStatusLabel: Record<Task["status"], string> = {
  todo: "К выполнению",
  in_progress: "В работе",
  done: "Готово",
};

const roleLabel: Record<Participant["role"], string> = {
  participant: "Участник",
  responsible: "Ответственный",
};

const rsvpLabel: Record<RSVPStatus, string> = {
  pending: "Ожидает",
  accepted: "Принял",
  declined: "Отклонил",
};

const rsvpClass: Record<RSVPStatus, string> = {
  pending: "rsvp--pending",
  accepted: "rsvp--accepted",
  declined: "rsvp--declined",
};

export function EventDetailPage() {
  const { id } = useParams<{ id: string }>();
  const eventId = Number(id);
  const { userId, role } = useAuth();

  const [event, setEvent] = useState<Event | null>(null);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [participants, setParticipants] = useState<Participant[]>([]);
  const [agenda, setAgenda] = useState<AgendaItem[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  const [newTaskTitle, setNewTaskTitle] = useState("");
  const [newTaskAssignee, setNewTaskAssignee] = useState("");
  const [newTaskDueDate, setNewTaskDueDate] = useState("");
  const [taskError, setTaskError] = useState("");

  const [newParticipantUserId, setNewParticipantUserId] = useState("");
  const [newParticipantRole, setNewParticipantRole] =
    useState<Participant["role"]>("participant");
  const [participantError, setParticipantError] = useState("");

  const [newAgendaTitle, setNewAgendaTitle] = useState("");
  const [newAgendaDuration, setNewAgendaDuration] = useState("");
  const [agendaError, setAgendaError] = useState("");

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
        setError(friendlyError(err, "Не удалось загрузить встречу")),
      )
      .finally(() => setLoading(false));
  }, [eventId]);

  async function handleCreateTask(e: React.FormEvent) {
    e.preventDefault();
    setTaskError("");

    try {
      const created = await createTask(eventId, {
        title: newTaskTitle,
        description: "",
        status: "todo",
        assignee_id: newTaskAssignee ? Number(newTaskAssignee) : null,
        due_date: newTaskDueDate
          ? new Date(newTaskDueDate).toISOString()
          : null,
      });
      setTasks((prev) => [...prev, created]);
      setNewTaskTitle("");
      setNewTaskAssignee("");
      setNewTaskDueDate("");
    } catch (err) {
      setTaskError(friendlyError(err, "Не удалось создать задачу"));
    }
  }

  async function handleChangeStatus(task: Task, status: Task["status"]) {
    setTaskError("");
    try {
      const updated = await updateTask(eventId, task.id, {
        title: task.title,
        description: task.description,
        status,
        assignee_id: task.assignee_id,
        due_date: task.due_date,
      });
      setTasks((prev) => prev.map((t) => (t.id === task.id ? updated : t)));
    } catch (err) {
      setTaskError(friendlyError(err, "Не удалось обновить задачу"));
    }
  }

  async function handleDeleteTask(taskId: number) {
    setTaskError("");
    try {
      await deleteTask(eventId, taskId);
      setTasks((prev) => prev.filter((t) => t.id !== taskId));
    } catch (err) {
      setTaskError(friendlyError(err, "Не удалось удалить задачу"));
    }
  }

  async function handleAddParticipant(e: React.FormEvent) {
    e.preventDefault();
    setParticipantError("");

    if (!newParticipantUserId) return;

    try {
      const created = await addParticipant(eventId, {
        user_id: Number(newParticipantUserId),
        role: newParticipantRole,
      });
      setParticipants((prev) => [...prev, created]);
      setNewParticipantUserId("");
      setNewParticipantRole("participant");
    } catch (err) {
      setParticipantError(friendlyError(err, "Не удалось добавить участника"));
    }
  }

  async function handleChangeParticipantRole(
    participant: Participant,
    newRole: Participant["role"],
  ) {
    setParticipantError("");
    try {
      const updated = await updateParticipant(eventId, participant.id, { role: newRole });
      setParticipants((prev) =>
        prev.map((p) => (p.id === participant.id ? updated : p)),
      );
    } catch (err) {
      setParticipantError(friendlyError(err, "Не удалось обновить роль"));
    }
  }

  async function handleRSVP(participantId: number, rsvpStatus: RSVPStatus) {
    setParticipantError("");
    try {
      const updated = await updateRSVP(eventId, participantId, { rsvp_status: rsvpStatus });
      setParticipants((prev) =>
        prev.map((p) => (p.id === participantId ? updated : p)),
      );
    } catch (err) {
      setParticipantError(friendlyError(err, "Не удалось обновить статус"));
    }
  }

  async function handleRemoveParticipant(participantId: number) {
    setParticipantError("");
    try {
      await removeParticipant(eventId, participantId);
      setParticipants((prev) => prev.filter((p) => p.id !== participantId));
    } catch (err) {
      setParticipantError(friendlyError(err, "Не удалось удалить участника"));
    }
  }

  async function handleCreateAgenda(e: React.FormEvent) {
    e.preventDefault();
    setAgendaError("");

    try {
      const created = await createAgendaItem(eventId, {
        title: newAgendaTitle,
        description: "",
        duration_minutes: newAgendaDuration ? Number(newAgendaDuration) : null,
      });
      setAgenda((prev) => [...prev, created]);
      setNewAgendaTitle("");
      setNewAgendaDuration("");
    } catch (err) {
      setAgendaError(friendlyError(err, "Не удалось добавить пункт"));
    }
  }

  async function handleToggleAgendaDone(item: AgendaItem) {
    setAgendaError("");
    try {
      const updated = await updateAgendaItem(eventId, item.id, {
        title: item.title,
        description: item.description,
        duration_minutes: item.duration_minutes,
        is_done: !item.is_done,
      });
      setAgenda((prev) => prev.map((a) => (a.id === item.id ? updated : a)));
    } catch (err) {
      setAgendaError(friendlyError(err, "Не удалось обновить пункт"));
    }
  }

  async function handleDeleteAgenda(itemId: number) {
    setAgendaError("");
    try {
      await deleteAgendaItem(eventId, itemId);
      setAgenda((prev) => prev.filter((a) => a.id !== itemId));
    } catch (err) {
      setAgendaError(friendlyError(err, "Не удалось удалить пункт"));
    }
  }

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
  const participantUserIds = new Set(participants.map((p) => p.user_id));
  const availableUsers = users.filter((u) => !participantUserIds.has(u.id));

  const isAdmin = role === "admin";
  const isCreator = userId != null && event.creator_id === userId;
  const canManageEvent = isAdmin || isCreator;

  const myParticipation = participants.find((p) => p.user_id === userId);

  function canEditTask(task: Task) {
    return (
      isAdmin ||
      isCreator ||
      (task.assignee_id != null && task.assignee_id === userId)
    );
  }

  return (
    <div className="event-detail-page">
      <Link to="/events" className="back-link">
        ← Все мероприятия
      </Link>

      <header className="event-detail-header">
        <div className="event-detail-header__title-row">
          <h1>{event.title}</h1>
          <span className={`status-badge ${statusClass[event.status]}`}>
            {statusLabel[event.status]}
          </span>
        </div>
        <p className="event-detail-time">
          {new Date(event.start_time).toLocaleString()} —{" "}
          {new Date(event.end_time).toLocaleString()}
        </p>
        {event.location && (
          <p className="event-detail-meta">📍 {event.location}</p>
        )}
        {event.meeting_url && (
          <p className="event-detail-meta">
            🔗{" "}
            <a href={event.meeting_url} target="_blank" rel="noopener noreferrer">
              Ссылка на встречу
            </a>
          </p>
        )}
        <p className="event-detail-description">
          {event.description || "Без описания"}
        </p>
        <p className="event-detail-creator">
          Создал: {userById.get(event.creator_id)?.email ?? `#${event.creator_id}`}
        </p>

        {myParticipation && (
          <div className="rsvp-block">
            <span className="rsvp-block__label">
              Мой статус:{" "}
              <span className={`rsvp-badge ${rsvpClass[myParticipation.rsvp_status]}`}>
                {rsvpLabel[myParticipation.rsvp_status]}
              </span>
            </span>
            <div className="rsvp-block__actions">
              {myParticipation.rsvp_status !== "accepted" && (
                <button
                  className="rsvp-btn rsvp-btn--accept"
                  onClick={() => handleRSVP(myParticipation.id, "accepted")}
                >
                  Принять
                </button>
              )}
              {myParticipation.rsvp_status !== "declined" && (
                <button
                  className="rsvp-btn rsvp-btn--decline"
                  onClick={() => handleRSVP(myParticipation.id, "declined")}
                >
                  Отклонить
                </button>
              )}
            </div>
          </div>
        )}
      </header>

      <div className="event-detail-grid">
        <section className="panel">
          <h2>Задачи</h2>

          {tasks.length === 0 ? (
            <p className="empty-text">Задач пока нет</p>
          ) : (
            <ul className="task-list">
              {tasks.map((t) => {
                const editable = canEditTask(t);
                return (
                  <li key={t.id} className={`task-item task-status-${t.status}`}>
                    <div className="task-row">
                      <div className="task-row__main">
                        <div className="task-item__title">{t.title}</div>
                        <div className="task-item__meta">
                          {editable ? (
                            <select
                              className="task-row__status"
                              value={t.status}
                              onChange={(e) =>
                                handleChangeStatus(t, e.target.value as Task["status"])
                              }
                            >
                              <option value="todo">{taskStatusLabel.todo}</option>
                              <option value="in_progress">{taskStatusLabel.in_progress}</option>
                              <option value="done">{taskStatusLabel.done}</option>
                            </select>
                          ) : (
                            <span className="task-row__status task-row__status--readonly">
                              {taskStatusLabel[t.status]}
                            </span>
                          )}
                          {t.assignee_id != null && (
                            <span>
                              {userById.get(t.assignee_id)?.email ?? `#${t.assignee_id}`}
                            </span>
                          )}
                          {t.due_date && (
                            <span>до {new Date(t.due_date).toLocaleDateString()}</span>
                          )}
                        </div>
                      </div>
                      {(isAdmin || isCreator) && (
                        <button
                          className="task-row__delete"
                          onClick={() => handleDeleteTask(t.id)}
                          title="Удалить задачу"
                        >
                          ×
                        </button>
                      )}
                    </div>
                  </li>
                );
              })}
            </ul>
          )}

          {canManageEvent && (
            <form className="inline-form" onSubmit={handleCreateTask}>
              <input
                placeholder="Новая задача"
                value={newTaskTitle}
                onChange={(e) => setNewTaskTitle(e.target.value)}
                required
              />
              <select
                value={newTaskAssignee}
                onChange={(e) => setNewTaskAssignee(e.target.value)}
              >
                <option value="">Без исполнителя</option>
                {users.map((u) => (
                  <option key={u.id} value={u.id}>{u.email}</option>
                ))}
              </select>
              <input
                type="date"
                value={newTaskDueDate}
                onChange={(e) => setNewTaskDueDate(e.target.value)}
              />
              <button type="submit">Добавить</button>
            </form>
          )}
          {taskError && <p className="form-error">{taskError}</p>}
        </section>

        <section className="panel">
          <h2>Участники</h2>

          {participants.length === 0 ? (
            <p className="empty-text">Пока никого нет</p>
          ) : (
            <ul className="participant-list">
              {participants.map((p) => (
                <li key={p.id} className="participant-item">
                  <span className="participant-item__email">
                    {userById.get(p.user_id)?.email ?? `#${p.user_id}`}
                  </span>
                  <div className="participant-item__actions">
                    <span className={`rsvp-badge rsvp-badge--sm ${rsvpClass[p.rsvp_status]}`}>
                      {rsvpLabel[p.rsvp_status]}
                    </span>
                    {canManageEvent ? (
                      <select
                        className="participant-role-select"
                        value={p.role}
                        onChange={(e) =>
                          handleChangeParticipantRole(p, e.target.value as Participant["role"])
                        }
                      >
                        <option value="participant">{roleLabel.participant}</option>
                        <option value="responsible">{roleLabel.responsible}</option>
                      </select>
                    ) : (
                      <span className="participant-role">{roleLabel[p.role]}</span>
                    )}
                    {canManageEvent && (
                      <button
                        className="task-row__delete"
                        onClick={() => handleRemoveParticipant(p.id)}
                        title="Удалить участника"
                      >
                        ×
                      </button>
                    )}
                  </div>
                </li>
              ))}
            </ul>
          )}

          {canManageEvent && (
            <form className="inline-form" onSubmit={handleAddParticipant}>
              <select
                value={newParticipantUserId}
                onChange={(e) => setNewParticipantUserId(e.target.value)}
                required
              >
                <option value="">Выберите пользователя</option>
                {availableUsers.map((u) => (
                  <option key={u.id} value={u.id}>{u.email}</option>
                ))}
              </select>
              <select
                value={newParticipantRole}
                onChange={(e) =>
                  setNewParticipantRole(e.target.value as Participant["role"])
                }
              >
                <option value="participant">{roleLabel.participant}</option>
                <option value="responsible">{roleLabel.responsible}</option>
              </select>
              <button type="submit" disabled={availableUsers.length === 0}>
                Добавить
              </button>
            </form>
          )}
          {participantError && <p className="form-error">{participantError}</p>}
        </section>

        <section className="panel">
          <h2>Повестка</h2>
          {agenda.length === 0 ? (
            <p className="empty-text">Повестка пока пустая</p>
          ) : (
            <ul className="agenda-list">
              {agenda.map((a) => (
                <li
                  key={a.id}
                  className={`agenda-item ${a.is_done ? "agenda-item--done" : ""}`}
                >
                  {canManageEvent ? (
                    <button
                      className="agenda-item__toggle"
                      onClick={() => handleToggleAgendaDone(a)}
                      title={a.is_done ? "Отметить как не выполнено" : "Отметить как обсуждено"}
                    >
                      {a.is_done ? "✓" : a.position}
                    </button>
                  ) : (
                    <span className="agenda-item__toggle agenda-item__toggle--readonly">
                      {a.is_done ? "✓" : a.position}
                    </span>
                  )}
                  <div className="agenda-item__body">
                    <div className="agenda-item__title">{a.title}</div>
                    {a.description && (
                      <div className="agenda-item__description">{a.description}</div>
                    )}
                    {a.duration_minutes != null && (
                      <div className="agenda-item__duration">≈ {a.duration_minutes} мин</div>
                    )}
                  </div>
                  {canManageEvent && (
                    <button
                      className="task-row__delete"
                      onClick={() => handleDeleteAgenda(a.id)}
                      title="Удалить пункт"
                    >
                      ×
                    </button>
                  )}
                </li>
              ))}
            </ul>
          )}

          {canManageEvent && (
            <form className="inline-form" onSubmit={handleCreateAgenda}>
              <input
                placeholder="Новый пункт повестки"
                value={newAgendaTitle}
                onChange={(e) => setNewAgendaTitle(e.target.value)}
                required
              />
              <input
                type="number"
                placeholder="Длительность, мин"
                value={newAgendaDuration}
                onChange={(e) => setNewAgendaDuration(e.target.value)}
                min={0}
              />
              <button type="submit">Добавить</button>
            </form>
          )}
          {agendaError && <p className="form-error">{agendaError}</p>}
        </section>
      </div>
    </div>
  );
}
