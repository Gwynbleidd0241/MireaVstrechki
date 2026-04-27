import { apiRequest } from "./client";

export type Event = {
  id: number;
  title: string;
  description: string;
  start_time: string;
  end_time: string;
  creator_id: number;
  created_at: string;
};

export type CreateEventRequest = {
  title: string;
  description: string;
  start_time: string;
  end_time: string;
};

export function getEvents() {
  return apiRequest<Event[]>("/events");
}

export function createEvent(data: CreateEventRequest) {
  return apiRequest<Event>("/events", {
    method: "POST",
    body: JSON.stringify(data),
  });
}
