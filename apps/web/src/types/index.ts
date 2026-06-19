export type Role = "admin" | "user";

export interface User {
  id: string;
  username: string;
  role: Role;
}

export interface Announcement {
  id: string;
  title: string;
  content: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface AnnouncementWithStatus {
  id: string;
  title: string;
  content: string;
  created_by: string;
  created_at: string;
  updated_at: string;
  is_read: boolean;
}

export type AnnouncementEventType = "created" | "updated" | "deleted";

export interface AnnouncementEvent {
  type: AnnouncementEventType;
  announcement?: Announcement;
  id?: string;
}

export interface AuthResponse {
  user: User;
}

export interface UnreadCountResponse {
  unread_count: number;
}

export type ToastType = "success" | "error" | "info";

export interface ToastItem {
  id: number;
  message: string;
  type: ToastType;
  title?: string;
  duration?: number;
}
