import { apiRequest } from "./client";

export type AuthResponse = {
  id: number;
  email: string;
  role: string;
  token?: string;
};

export function register(email: string, password: string, role: string) {
  return apiRequest<AuthResponse>("/register", {
    method: "POST",
    body: JSON.stringify({ email, password, role }),
  });
}

export function login(email: string, password: string) {
  return apiRequest<AuthResponse>("/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}
