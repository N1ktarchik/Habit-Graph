import {
  ActivityGraph,
  CreateHabitDTO,
  Habit,
  TrackActivityDTO,
  UpdateHabitDTO,
} from '../types/habit';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    headers: {
      'Content-Type': 'application/json',
    },
    ...options,
  });

  if (!response.ok) {
    let errorMessage = `Ошибка запроса (${response.status})`;
    try {
      const errData = await response.json();
      if (errData?.message) {
        errorMessage = errData.message;
      }
    } catch {
      // Игнорируем ошибку парсинга JSON, если ответ пустой или текстовый
    }
    throw new Error(errorMessage);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
}

export const habitsApi = {
  getHabits: () => request<Habit[]>('/api/habits'),

  createHabit: (data: CreateHabitDTO) =>
    request<Habit>('/api/habits', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  updateHabit: (id: string, data: UpdateHabitDTO) =>
    request<Habit>(`/api/habits/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    }),

  deleteHabit: (id: string) =>
    request<void>(`/api/habits/${id}`, {
      method: 'DELETE',
    }),

  getGraph: (id: string, year: number) =>
    request<ActivityGraph>(`/api/habits/${id}/graph?year=${year}`),

  trackActivity: (id: string, data: TrackActivityDTO) =>
    request<{ status: string }>(`/api/habits/${id}/activities`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),
};