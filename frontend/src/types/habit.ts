export interface Habit {
  id: string;
  title: string;
  description: string;
  created_at: string;
}

export interface CreateHabitDTO {
  title: string;
  description?: string;
}

export interface UpdateHabitDTO {
  title?: string;
  description?: string;
}

export interface TrackActivityDTO {
  date: string; // YYYY-MM-DD
  value: number; // минуты
}

export type ActivityGraph = Record<string, number>;

export interface ApiError {
  code: number;
  message: string;
}