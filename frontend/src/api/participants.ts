import { apiRequest } from "./client";

export type ParticipantRole = "participant" | "responsible";

export type Participant = {
  id: number;
  event_id: number;
  user_id: number;
  role: ParticipantRole;
  created_at: string;
};

export type AddParticipantRequest = {
  user_id: number;
  role: ParticipantRole;
};

export type UpdateParticipantRequest = {
  role: ParticipantRole;
};

export function getParticipants(eventId: number) {
  return apiRequest<Participant[]>(`/events/${eventId}/participants`);
}

export function addParticipant(
  eventId: number,
  data: AddParticipantRequest,
) {
  return apiRequest<Participant>(`/events/${eventId}/participants`, {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export function updateParticipant(
  eventId: number,
  participantId: number,
  data: UpdateParticipantRequest,
) {
  return apiRequest<Participant>(
    `/events/${eventId}/participants/${participantId}`,
    {
      method: "PATCH",
      body: JSON.stringify(data),
    },
  );
}

export function removeParticipant(eventId: number, participantId: number) {
  return apiRequest<void>(
    `/events/${eventId}/participants/${participantId}`,
    {
      method: "DELETE",
    },
  );
}
