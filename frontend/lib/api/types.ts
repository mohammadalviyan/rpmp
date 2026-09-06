export type UserRole = "viewer" | "admin";

export type User = {
  id: string;
  employee_id: string;
  display_name: string;
  role: UserRole;
};

export type AuthResponse = {
  user: User;
};

export type LoginRequest = {
  employee_id: string;
  password: string;
};

export type ApiError = {
  code: string;
  message: string;
};
