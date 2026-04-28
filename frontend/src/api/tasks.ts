import { apiRequest } from "./client";

export type TaskStatus = "todo" | "in_progress" | "done";

export type Task = {
  id: number;
  event_id: number;
  title: string;
  description: string;
  status: TaskStatus;
  assignee_id: number | null;
  due_date: string | null;
  created_at: string;
};

export type CreateTaskRequest = {
  title: string;
  description: string;
  status: TaskStatus;
  assignee_id: number | null;
  due_date: string | null;
};

export type UpdateTaskRequest = CreateTaskRequest;

export function getTasks(eventId: number) {
  return apiRequest<Task[]>(`/events/${eventId}/tasks`);
}

export function createTask(eventId: number, data: CreateTaskRequest) {
  return apiRequest<Task>(`/events/${eventId}/tasks`, {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export function updateTask(
  eventId: number,
  taskId: number,
  data: UpdateTaskRequest,
) {
  return apiRequest<Task>(`/events/${eventId}/tasks/${taskId}`, {
    method: "PATCH",
    body: JSON.stringify(data),
  });
}

export function deleteTask(eventId: number, taskId: number) {
  return apiRequest<void>(`/events/${eventId}/tasks/${taskId}`, {
    method: "DELETE",
  });
}
