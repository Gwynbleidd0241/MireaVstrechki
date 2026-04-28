import { apiRequest } from "./client";

export type User = {
  id: number;
  email: string;
  role: string;
};

export function getUsers() {
  return apiRequest<User[]>("/users");
}
