import React, { useEffect, useState } from 'react';
import {
  Play,
  Square,
  Pause,
  RotateCcw,
  ChevronDown,
  ChevronUp,
  Clock,
  Calendar,
  Sparkles,
} from 'lucide-react';
import { formatTimerSeconds } from '../utils/date';

interface PersistedTimerState {
  habitId: string;
  isRunning: boolean;
  startTimestamp: number | null; // ms
  accumulatedSeconds: number;
}

interface ActivityTimerProps {
  habitId: string;
  habitTitle: string;
  selectedDate: string;
  onSelectDate: (date: string) => void;
  onTrack: (date: string, minutes: number) => Promise<void>;
}

const STORAGE_KEY = 'habit_graph_active_timer_v1';

export const ActivityTimer: React.FC<ActivityTimerProps> = ({
  habitId,
  habitTitle,
  selectedDate,
  onSelectDate,
  onTrack,
}) => {
  const [timerState, setTimerState] = useState<PersistedTimerState>(() => {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (raw) {
        const parsed = JSON.parse(raw) as PersistedTimerState;
        if (parsed && typeof parsed.habitId === 'string') {
          return parsed;
        }
      }
    } catch {
      // ignore localStorage errors
    }
    return {
      habitId,
      isRunning: false,
      startTimestamp: null,
      accumulatedSeconds: 0,
    };
  });

  const [displaySeconds, setDisplaySeconds] = useState<number>(0);
  const [submitting, setSubmitting] = useState<boolean>(false);

  // Состояние для ручного режима (аккордеон)
  const [showManual, setShowManual] = useState<boolean>(false);
  const [manualMinutes, setManualMinutes] = useState<number>(30);

  // Синхронизация с localStorage
  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(timerState));
    } catch {
      // ignore
    }
  }, [timerState]);

  // Вычисление текущего времени таймера каждую секунду
  useEffect(() => {
    const computeElapsed = () => {
      if (timerState.habitId !== habitId) return 0;
      let total = timerState.accumulatedSeconds;
      if (timerState.isRunning && timerState.startTimestamp !== null) {
        total += Math.floor((Date.now() - timerState.startTimestamp) / 1000);
      }
      return Math.max(0, total);
    };

    setDisplaySeconds(computeElapsed());

    if (!timerState.isRunning || timerState.habitId !== habitId) return;

    const interval = window.setInterval(() => {
      setDisplaySeconds(computeElapsed());
    }, 500);

    return () => window.clearInterval(interval);
  }, [timerState, habitId]);

  const isTimerForAnotherHabit =
    timerState.habitId !== habitId &&
    (timerState.isRunning || timerState.accumulatedSeconds > 0);

  const isRunningCurrent = timerState.habitId === habitId && timerState.isRunning;
  const hasElapsedCurrent = timerState.habitId === habitId && displaySeconds > 0;

  const handleStart = () => {
    setTimerState({
      habitId,
      isRunning: true,
      startTimestamp: Date.now(),
      accumulatedSeconds: timerState.habitId === habitId ? timerState.accumulatedSeconds : 0,
    });
  };

  const handlePause = () => {
    if (!isRunningCurrent || timerState.startTimestamp === null) return;
    const delta = Math.floor((Date.now() - timerState.startTimestamp) / 1000);
    setTimerState({
      habitId,
      isRunning: false,
      startTimestamp: null,
      accumulatedSeconds: timerState.accumulatedSeconds + Math.max(0, delta),
    });
  };

  const handleReset = () => {
    setTimerState({
      habitId,
      isRunning: false,
      startTimestamp: null,
      accumulatedSeconds: 0,
    });
    setDisplaySeconds(0);
  };

  const handleFinish = async () => {
    if (displaySeconds <= 0 || submitting) return;
    setSubmitting(true);

    // Округляем до минут, минимум 1 минута при любом запуске
    const calculatedMinutes = Math.max(1, Math.round(displaySeconds / 60));
    const todayStr = new Date().toISOString().split('T')[0];

    try {
      await onTrack(todayStr, calculatedMinutes);
      onSelectDate(todayStr);
      handleReset();
    } finally {
      setSubmitting(false);
    }
  };

  const handleManualSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (submitting) return;
    setSubmitting(true);
    try {
      await onTrack(selectedDate, Number(manualMinutes));
      setShowManual(false);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="rounded-xl bg-zinc-900/60 border border-zinc-800/80 overflow-hidden">
      {/* Основной блок таймера */}
      <div className="p-6 flex flex-col sm:flex-row items-center justify-between gap-6">
        <div className="space-y-1.5 text-center sm:text-left">
          <div className="inline-flex items-center gap-2 text-xs font-medium text-emerald-400 bg-emerald-500/10 border border-emerald-500/20 px-2.5 py-1 rounded-full">
            <Sparkles className="h-3.5 w-3.5" />
            <span>Трекер реального времени</span>
          </div>
          <h3 className="text-base font-semibold text-zinc-100">
            Сессия: <span className="text-emerald-400">{habitTitle}</span>
          </h3>
          <p className="text-xs text-zinc-400">
            Нажми «Начать занятие», сфокусируйся на деле, а по окончании нажми «Завершить» — время запишется автоматически.
          </p>
          {isTimerForAnotherHabit && (
            <p className="text-xs text-amber-400 pt-1">
              ⚠️ Внимание: сейчас активен таймер для другой привычки. Запуск здесь сбросит его.
            </p>
          )}
        </div>

        {/* Циферблат и кнопки */}
        <div className="flex flex-col items-center gap-3 shrink-0">
          <div
            className={`font-mono text-3xl sm:text-4xl font-bold tracking-wider px-5 py-2.5 rounded-xl border transition-all ${
              isRunningCurrent
                ? 'bg-emerald-950/40 border-emerald-500/50 text-emerald-300 shadow-[0_0_25px_rgba(16,185,129,0.15)]'
                : hasElapsedCurrent
                ? 'bg-zinc-900 border-amber-500/40 text-amber-300'
                : 'bg-zinc-950 border-zinc-800 text-zinc-300'
            }`}
          >
            {formatTimerSeconds(displaySeconds)}
          </div>

          <div className="flex items-center gap-2">
            {!isRunningCurrent ? (
              <button
                type="button"
                onClick={handleStart}
                disabled={submitting}
                className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-emerald-500 hover:bg-emerald-400 text-zinc-950 font-semibold text-sm transition shadow-[0_0_15px_rgba(16,185,129,0.25)]"
              >
                <Play className="h-4 w-4 fill-current" />
                {hasElapsedCurrent ? 'Продолжить' : 'Начать занятие'}
              </button>
            ) : (
              <button
                type="button"
                onClick={handlePause}
                disabled={submitting}
                className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-amber-500/20 border border-amber-500/40 hover:bg-amber-500/30 text-amber-300 font-medium text-sm transition"
              >
                <Pause className="h-4 w-4" />
                Пауза
              </button>
            )}

            {hasElapsedCurrent && (
              <>
                <button
                  type="button"
                  onClick={handleFinish}
                  disabled={submitting}
                  className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-zinc-100 hover:bg-white text-zinc-950 font-semibold text-sm transition"
                >
                  <Square className="h-4 w-4 fill-current" />
                  Завершить ({Math.max(1, Math.round(displaySeconds / 60))} м)
                </button>
                <button
                  type="button"
                  onClick={handleReset}
                  disabled={submitting}
                  title="Сбросить таймер"
                  className="p-2 rounded-lg bg-zinc-950 border border-zinc-800 text-zinc-400 hover:text-red-400 hover:border-red-500/40 transition"
                >
                  <RotateCcw className="h-4 w-4" />
                </button>
              </>
            )}
          </div>
        </div>
      </div>

      {/* Аккордеон для ручного ввода */}
      <div className="border-t border-zinc-800/80 bg-zinc-950/40">
        <button
          type="button"
          onClick={() => setShowManual((prev) => !prev)}
          className="w-full px-6 py-2.5 flex items-center justify-between text-xs text-zinc-400 hover:text-zinc-200 transition"
        >
          <span>Забыл включить таймер? Добавить или изменить время вручную</span>
          {showManual ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
        </button>

        {showManual && (
          <form
            onSubmit={handleManualSubmit}
            className="px-6 pb-5 pt-2 flex flex-wrap items-end gap-4 border-t border-zinc-800/40"
          >
            <div className="flex-1 min-w-[180px] space-y-1.5">
              <label className="text-xs text-zinc-400 flex items-center gap-1.5">
                <Calendar className="h-3.5 w-3.5 text-emerald-400" />
                Дата (YYYY-MM-DD)
              </label>
              <input
                type="date"
                required
                min="2026-01-01"
                value={selectedDate}
                onChange={(e) => onSelectDate(e.target.value)}
                className="w-full px-3.5 py-2 rounded-lg bg-zinc-950 border border-zinc-800 text-sm focus:outline-none focus:border-emerald-500"
              />
            </div>

            <div className="flex-1 min-w-[180px] space-y-1.5">
              <label className="text-xs text-zinc-400 flex items-center gap-1.5">
                <Clock className="h-3.5 w-3.5 text-emerald-400" />
                Время за день (минуты)
              </label>
              <input
                type="number"
                required
                min={0}
                max={1440}
                value={manualMinutes}
                onChange={(e) => setManualMinutes(Number(e.target.value))}
                className="w-full px-3.5 py-2 rounded-lg bg-zinc-950 border border-zinc-800 text-sm focus:outline-none focus:border-emerald-500"
              />
            </div>

            <button
              type="submit"
              disabled={submitting}
              className="px-5 py-2 rounded-lg bg-zinc-800 hover:bg-zinc-700 text-zinc-100 font-medium text-sm transition"
            >
              Сохранить вручную
            </button>
          </form>
        )}
      </div>
    </div>
  );
};