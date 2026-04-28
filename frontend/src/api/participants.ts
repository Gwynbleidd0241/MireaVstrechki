import { apiRequest } from "./client";

export type ParticipantRole = "participant" | "responsible";

export type Participant = {
  id: number;
  event_id: number;
  user_id: number;
  role: ParticipantRole;
  created_at: string;
};

export function getParticipants(eventId: number) {
  return apiRequest<Participant[]>(`/events/${eventId}/participants`);
}
