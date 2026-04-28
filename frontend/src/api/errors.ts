const messages: Record<string, string> = {
  "permission denied": "Недостаточно прав",
  "unauthorized": "Сессия истекла, войдите заново",
  "invalid json": "Ошибка формата данных",
  "invalid email": "Введите корректный email",
  "password too short": "Пароль должен быть не менее 8 символов",
  "password too long": "Пароль слишком длинный",
  "invalid role": "Неверная роль",
  "invalid email or password": "Неверный email или пароль",
  "event title required": "Введите название встречи",
  "event title too long": "Название встречи слишком длинное",
  "event description too long": "Описание слишком длинное",
  "invalid event time": "Время окончания должно быть позже начала",
  "event not found": "Встреча не найдена",
  "task title required": "Введите название задачи",
  "task title too long": "Название задачи слишком длинное",
  "task description too long": "Описание задачи слишком длинное",
  "invalid task status": "Недопустимый статус задачи",
  "task not found": "Задача не найдена",
  "participant user required": "Выберите пользователя",
  "invalid participant role": "Неверная роль участника",
  "participant not found": "Участник не найден",
  "agenda item title required": "Введите название пункта",
  "agenda item title too long": "Название пункта слишком длинное",
  "agenda item description too long": "Описание пункта слишком длинное",
  "invalid agenda item duration": "Длительность не может быть отрицательной",
  "agenda item not found": "Пункт повестки не найден",
  "invalid event id": "Неверный идентификатор встречи",
  "invalid path": "Некорректный путь запроса",
  "invalid start_time": "Некорректное время начала",
  "invalid end_time": "Некорректное время окончания",
  "invalid due_date": "Некорректный дедлайн",
};

export function friendlyError(err: unknown, fallback = "Что-то пошло не так"): string {
  if (err instanceof Error) {
    const key = err.message.trim().toLowerCase();
    return messages[key] || fallback;
  }
  return fallback;
}
