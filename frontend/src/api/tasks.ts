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

export function getTasks(eventId: number) {
  return apiRequest<Task[]>(`/events/${eventId}/tasks`);
}
