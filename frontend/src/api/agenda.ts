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

export type CreateAgendaItemRequest = {
  title: string;
  description: string;
  duration_minutes: number | null;
};

export type UpdateAgendaItemRequest = {
  title: string;
  description: string;
  duration_minutes: number | null;
  is_done: boolean;
};

export function getAgenda(eventId: number) {
  return apiRequest<AgendaItem[]>(`/events/${eventId}/agenda`);
}

export function createAgendaItem(
  eventId: number,
  data: CreateAgendaItemRequest,
) {
  return apiRequest<AgendaItem>(`/events/${eventId}/agenda`, {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export function updateAgendaItem(
  eventId: number,
  itemId: number,
  data: UpdateAgendaItemRequest,
) {
  return apiRequest<AgendaItem>(`/events/${eventId}/agenda/${itemId}`, {
    method: "PATCH",
    body: JSON.stringify(data),
  });
}

export function deleteAgendaItem(eventId: number, itemId: number) {
  return apiRequest<void>(`/events/${eventId}/agenda/${itemId}`, {
    method: "DELETE",
  });
}
