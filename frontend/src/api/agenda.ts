import { apiRequest } from "./client";

export type AgendaItem = {
  id: number;
  event_id: number;
  position: number;
  title: string;
  description: string;
  duration_minutes: number | null;
  is_done: boolean;
  created_at: string;
};

export function getAgenda(eventId: number) {
  return apiRequest<AgendaItem[]>(`/events/${eventId}/agenda`);
}
