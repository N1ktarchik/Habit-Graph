import React, { useEffect, useState } from 'react';
import {
  Activity,
  Calendar,
  Edit3,
  Plus,
  RefreshCw,
  Trash2,
  Check,
  X,
  LayoutGrid,
  CalendarDays,
  Info,
} from 'lucide-react';
import { habitsApi } from './api/client';
import { ActivityGraph, Habit } from './types/habit';
import { Heatmap } from './components/Heatmap';
import { MonthHeatmap } from './components/MonthHeatmap';
import { MonthTrendChart } from './components/MonthTrendChart';
import { ActivityTimer } from './components/ActivityTimer';
import {
  formatDuration,
  formatHumanDate,
  getIntensityInfo,
  MONTH_NAMES,
} from './utils/date';

type ViewMode = 'month' | 'year';

export const App: React.FC = () => {
  const now = new Date();
  const currentYear = Math.max(2026, now.getFullYear());

  const [habits, setHabits] = useState<Habit[]>([]);
  const [selectedHabitId, setSelectedHabitId] = useState<string | null>(null);
  const [graph, setGraph] = useState<ActivityGraph>({});

  // Период и режим отображения
  const [viewMode, setViewMode] = useState<ViewMode>('month');
  const [year, setYear] = useState<number>(currentYear);
  const [month, setMonth] = useState<number>(() => now.getMonth());

  // Выбранная дата (для просмотра деталей дня и ручного редактирования)
  const [selectedDate, setSelectedDate] = useState<string>(
    () => now.toISOString().split('T')[0]
  );

  // Создание привычки
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');

  // Редактирование привычки
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editTitle, setEditTitle] = useState('');
  const [editDesc, setEditDesc] = useState('');

  // Статусы
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const selectedHabit = habits.find((h) => h.id === selectedHabitId) || null;

  const fetchHabits = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await habitsApi.getHabits();
      setHabits(data);
      if (data.length > 0 && !selectedHabitId) {
        setSelectedHabitId(data[0].id);
      }
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : 'Ошибка при загрузке привычек';
      setError(message);
    } finally {
      setLoading(false);
    }
  };

  const fetchGraph = async (habitId: string, targetYear: number) => {
    try {
      const data = await habitsApi.getGraph(habitId, targetYear);
      setGraph(data || {});
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : 'Ошибка при загрузке графика активности';
      setError(message);
    }
  };

  useEffect(() => {
    fetchHabits();
  }, []);

  useEffect(() => {
    if (selectedHabitId) {
      fetchGraph(selectedHabitId, year);
    } else {
      setGraph({});
    }
  }, [selectedHabitId, year]);

  const handleCreateHabit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      const created = await habitsApi.createHabit({
        title: title.trim(),
        description: description.trim(),
      });
      setHabits((prev) => [...prev, created]);
      setSelectedHabitId(created.id);
      setTitle('');
      setDescription('');
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Не удалось создать привычку');
    }
  };

  const handleUpdateHabit = async (id: string) => {
    setError(null);
    try {
      const updated = await habitsApi.updateHabit(id, {
        title: editTitle.trim(),
        description: editDesc.trim(),
      });
      setHabits((prev) => prev.map((h) => (h.id === id ? updated : h)));
      setEditingId(null);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Не удалось обновить привычку');
    }
  };

  const handleDeleteHabit = async (id: string) => {
    setError(null);
    try {
      await habitsApi.deleteHabit(id);
      const remaining = habits.filter((h) => h.id !== id);
      setHabits(remaining);
      if (selectedHabitId === id) {
        setSelectedHabitId(remaining.length > 0 ? remaining[0].id : null);
      }
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Не удалось удалить привычку');
    }
  };

  const handleTrackActivity = async (dateStr: string, minutesToSet: number) => {
    if (!selectedHabitId) return;
    setError(null);
    try {
      await habitsApi.trackActivity(selectedHabitId, {
        date: dateStr,
        value: minutesToSet,
      });
      await fetchGraph(selectedHabitId, year);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Не удалось сохранить активность');
    }
  };

  const selectedDayMinutes = graph[selectedDate] || 0;
  const selectedDayIntensity = getIntensityInfo(selectedDayMinutes);

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100">
      {/* Header */}
      <header className="border-b border-zinc-800/80 bg-zinc-950/80 backdrop-blur sticky top-0 z-30">
        <div className="max-w-6xl mx-auto px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-lg bg-emerald-500/10 border border-emerald-500/30 flex items-center justify-center text-emerald-400">
              <Activity className="h-5 w-5" />
            </div>
            <div>
              <h1 className="text-lg font-bold tracking-tight">Habit Graph</h1>
              <p className="text-xs text-zinc-400">Интерактивный трекер привычек</p>
            </div>
          </div>
          <button
            onClick={fetchHabits}
            className="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg bg-zinc-900 border border-zinc-800 text-xs text-zinc-300 hover:border-zinc-700 transition"
          >
            <RefreshCw className={`h-3.5 w-3.5 ${loading ? 'animate-spin' : ''}`} />
            Обновить
          </button>
        </div>
      </header>

      <main className="max-w-6xl mx-auto px-6 py-8 grid grid-cols-1 lg:grid-cols-12 gap-8">
        {/* Левая колонка: список и создание */}
        <section className="lg:col-span-4 space-y-6">
          <form
            onSubmit={handleCreateHabit}
            className="p-5 rounded-xl bg-zinc-900/60 border border-zinc-800/80 space-y-4"
          >
            <h2 className="text-sm font-semibold text-zinc-200 flex items-center gap-2">
              <Plus className="h-4 w-4 text-emerald-400" />
              Новая привычка
            </h2>
            <div className="space-y-3">
              <input
                type="text"
                required
                minLength={3}
                maxLength={120}
                placeholder="Название (до 120 символов)"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="w-full px-3.5 py-2 rounded-lg bg-zinc-950 border border-zinc-800 text-sm focus:outline-none focus:border-emerald-500"
              />
              <textarea
                rows={2}
                maxLength={200}
                placeholder="Описание (до 200 символов)"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                className="w-full px-3.5 py-2 rounded-lg bg-zinc-950 border border-zinc-800 text-sm focus:outline-none focus:border-emerald-500 resize-none"
              />
              <button
                type="submit"
                className="w-full py-2 px-4 rounded-lg bg-emerald-500 hover:bg-emerald-400 text-zinc-950 font-semibold text-sm transition"
              >
                Создать
              </button>
            </div>
          </form>

          {/* Список привычек */}
          <div className="space-y-2.5">
            <h3 className="text-xs font-semibold uppercase tracking-wider text-zinc-500 px-1">
              Мои привычки ({habits.length})
            </h3>
            {habits.length === 0 ? (
              <div className="p-6 rounded-xl border border-dashed border-zinc-800 text-center text-sm text-zinc-500">
                Список привычек пока пуст
              </div>
            ) : (
              habits.map((h) => {
                const isSelected = h.id === selectedHabitId;
                const isEditing = h.id === editingId;
                return (
                  <div
                    key={h.id}
                    onClick={() => !isEditing && setSelectedHabitId(h.id)}
                    className={`p-4 rounded-xl border transition cursor-pointer ${
                      isSelected
                        ? 'bg-zinc-900 border-emerald-500/50 shadow-[0_0_20px_rgba(16,185,129,0.08)]'
                        : 'bg-zinc-900/40 border-zinc-800/80 hover:border-zinc-700'
                    }`}
                  >
                    {isEditing ? (
                      <div className="space-y-2.5" onClick={(e) => e.stopPropagation()}>
                        <input
                          type="text"
                          value={editTitle}
                          onChange={(e) => setEditTitle(e.target.value)}
                          className="w-full px-2.5 py-1.5 rounded bg-zinc-950 border border-zinc-700 text-sm"
                        />
                        <input
                          type="text"
                          value={editDesc}
                          onChange={(e) => setEditDesc(e.target.value)}
                          className="w-full px-2.5 py-1.5 rounded bg-zinc-950 border border-zinc-700 text-xs"
                        />
                        <div className="flex justify-end gap-2">
                          <button
                            onClick={() => handleUpdateHabit(h.id)}
                            className="p-1.5 rounded bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30"
                          >
                            <Check className="h-4 w-4" />
                          </button>
                          <button
                            onClick={() => setEditingId(null)}
                            className="p-1.5 rounded bg-zinc-800 text-zinc-400 hover:bg-zinc-700"
                          >
                            <X className="h-4 w-4" />
                          </button>
                        </div>
                      </div>
                    ) : (
                      <div className="flex items-start justify-between gap-3">
                        <div className="space-y-1">
                          <h4 className="font-medium text-sm text-zinc-100">{h.title}</h4>
                          {h.description && (
                            <p className="text-xs text-zinc-400 line-clamp-2">{h.description}</p>
                          )}
                        </div>
                        <div
                          className="flex items-center gap-1 shrink-0"
                          onClick={(e) => e.stopPropagation()}
                        >
                          <button
                            onClick={() => {
                              setEditingId(h.id);
                              setEditTitle(h.title);
                              setEditDesc(h.description);
                            }}
                            className="p-1.5 rounded-lg text-zinc-400 hover:text-zinc-100 hover:bg-zinc-800"
                          >
                            <Edit3 className="h-3.5 w-3.5" />
                          </button>
                          <button
                            onClick={() => handleDeleteHabit(h.id)}
                            className="p-1.5 rounded-lg text-zinc-400 hover:text-red-400 hover:bg-red-500/10"
                          >
                            <Trash2 className="h-3.5 w-3.5" />
                          </button>
                        </div>
                      </div>
                    )}
                  </div>
                );
              })
            )}
          </div>
        </section>

        {/* Правая колонка: секундомер, хитмапы и аналитика месяца */}
        <section className="lg:col-span-8 space-y-6">
          {error && (
            <div className="p-4 rounded-xl bg-red-500/10 border border-red-500/30 text-red-300 text-sm flex items-center justify-between">
              <span>{error}</span>
              <button onClick={() => setError(null)} className="text-xs underline">
                Закрыть
              </button>
            </div>
          )}

          {selectedHabit ? (
            <>
              {/* 1. Интерактивный секундомер */}
              <ActivityTimer
                habitId={selectedHabit.id}
                habitTitle={selectedHabit.title}
                selectedDate={selectedDate}
                onSelectDate={(d) => {
                  setSelectedDate(d);
                  const parsedMonth = Number(d.split('-')[1]) - 1;
                  if (!Number.isNaN(parsedMonth) && parsedMonth >= 0 && parsedMonth <= 11) {
                    setMonth(parsedMonth);
                  }
                }}
                onTrack={handleTrackActivity}
              />

              {/* 2. Блок календаря активности (Месяц / Год) */}
              <div className="p-6 rounded-xl bg-zinc-900/60 border border-zinc-800/80 space-y-6">
                <div className="flex flex-wrap items-center justify-between gap-4">
                  <div>
                    <h2 className="text-xl font-bold text-zinc-100">{selectedHabit.title}</h2>
                    {selectedHabit.description && (
                      <p className="text-sm text-zinc-400 mt-0.5">{selectedHabit.description}</p>
                    )}
                  </div>

                  {/* Панель управления масштабом и датой */}
                  <div className="flex flex-wrap items-center gap-2">
                    {/* Переключатель Месяц / Год */}
                    <div className="inline-flex p-1 rounded-lg bg-zinc-950 border border-zinc-800">
                      <button
                        type="button"
                        onClick={() => setViewMode('month')}
                        className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-md text-xs font-medium transition ${
                          viewMode === 'month'
                            ? 'bg-emerald-500 text-zinc-950 font-semibold'
                            : 'text-zinc-400 hover:text-zinc-200'
                        }`}
                      >
                        <CalendarDays className="h-3.5 w-3.5" />
                        Месяц
                      </button>
                      <button
                        type="button"
                        onClick={() => setViewMode('year')}
                        className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-md text-xs font-medium transition ${
                          viewMode === 'year'
                            ? 'bg-emerald-500 text-zinc-950 font-semibold'
                            : 'text-zinc-400 hover:text-zinc-200'
                        }`}
                      >
                        <LayoutGrid className="h-3.5 w-3.5" />
                        Год
                      </button>
                    </div>

                    {/* Селектор месяца */}
                    <select
                      value={month}
                      onChange={(e) => setMonth(Number(e.target.value))}
                      className="bg-zinc-950 border border-zinc-800 rounded-lg px-3 py-1.5 text-xs sm:text-sm focus:outline-none focus:border-emerald-500"
                    >
                      {MONTH_NAMES.map((name, idx) => (
                        <option key={name} value={idx}>
                          {name}
                        </option>
                      ))}
                    </select>

                    {/* Селектор года */}
                    <div className="flex items-center gap-1.5">
                      <Calendar className="h-4 w-4 text-zinc-400 hidden sm:block" />
                      <select
                        value={year}
                        onChange={(e) => setYear(Number(e.target.value))}
                        className="bg-zinc-950 border border-zinc-800 rounded-lg px-3 py-1.5 text-xs sm:text-sm focus:outline-none focus:border-emerald-500"
                      >
                        <option value={2026}>2026</option>
                        <option value={2027}>2027</option>
                      </select>
                    </div>
                  </div>
                </div>

                {/* Отображение выбранного масштаба */}
                {viewMode === 'month' ? (
                  <MonthHeatmap
                    year={year}
                    month={month}
                    data={graph}
                    selectedDate={selectedDate}
                    onSelectDate={(dateStr) => setSelectedDate(dateStr)}
                  />
                ) : (
                  <Heatmap
                    year={year}
                    data={graph}
                    selectedDate={selectedDate}
                    onSelectDate={(dateStr) => {
                      setSelectedDate(dateStr);
                      const m = Number(dateStr.split('-')[1]) - 1;
                      if (!Number.isNaN(m) && m >= 0 && m <= 11) {
                        setMonth(m);
                      }
                    }}
                  />
                )}

                {/* Детали выбранного дня по клику */}
                <div className="pt-4 border-t border-zinc-800/80 flex flex-wrap items-center justify-between gap-3 text-xs">
                  <div className="flex items-center gap-2 text-zinc-300">
                    <Info className="h-4 w-4 text-emerald-400 shrink-0" />
                    <span>
                      Выбран день:{' '}
                      <strong className="text-zinc-100">
                        {formatHumanDate(selectedDate)}
                      </strong>
                    </span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-zinc-400">
                      Результат: <strong className="text-zinc-100">{formatDuration(selectedDayMinutes)}</strong>
                    </span>
                    <span
                      className={`px-2.5 py-0.5 rounded-full border text-[11px] font-medium ${selectedDayIntensity.badgeClass}`}
                    >
                      {selectedDayIntensity.label}
                    </span>
                  </div>
                </div>
              </div>

              {/* 3. График динамики за месяц */}
              <MonthTrendChart
                year={year}
                month={month}
                data={graph}
                selectedDate={selectedDate}
                onSelectDate={(dateStr) => setSelectedDate(dateStr)}
              />
            </>
          ) : (
            <div className="p-12 rounded-xl bg-zinc-900/30 border border-zinc-800/60 text-center space-y-2">
              <Activity className="h-8 w-8 text-zinc-600 mx-auto" />
              <p className="text-zinc-400 text-sm">
                Выбери привычку слева или создай новую, чтобы запустить таймер и увидеть график
              </p>
            </div>
          )}
        </section>
      </main>
    </div>
  );
};