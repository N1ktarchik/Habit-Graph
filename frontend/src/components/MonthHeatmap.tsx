import React, { useMemo, useState } from 'react';
import { ActivityGraph } from '../types/habit';
import {
  formatDuration,
  formatHumanDate,
  getIntensityInfo,
  MONTH_NAMES,
  SHORT_WEEKDAYS,
} from '../utils/date';

interface MonthHeatmapProps {
  year: number;
  month: number; // 0 .. 11
  data: ActivityGraph;
  selectedDate?: string;
  onSelectDate?: (dateStr: string) => void;
}

interface MonthDayItem {
  dayNumber: number;
  dateStr: string;
  value: number;
  dayOfWeek: number;
}

export const MonthHeatmap: React.FC<MonthHeatmapProps> = ({
  year,
  month,
  data,
  selectedDate,
  onSelectDate,
}) => {
  const [hoveredDay, setHoveredDay] = useState<{
    item: MonthDayItem;
    x: number;
    y: number;
  } | null>(null);

  const { cells, totalMinutes, activeDays } = useMemo(() => {
    const firstDay = new Date(Date.UTC(year, month, 1));
    const daysInMonth = new Date(Date.UTC(year, month + 1, 0)).getUTCDate();
    const startOffset = (firstDay.getUTCDay() + 6) % 7; // Пн = 0 .. Вс = 6

    const list: (MonthDayItem | null)[] = Array(startOffset).fill(null);
    let sum = 0;
    let active = 0;

    for (let d = 1; d <= daysInMonth; d++) {
      const dateObj = new Date(Date.UTC(year, month, d));
      const dateStr = dateObj.toISOString().split('T')[0];
      const val = data[dateStr] || 0;
      const dow = (dateObj.getUTCDay() + 6) % 7;

      sum += val;
      if (val > 0) active++;

      list.push({
        dayNumber: d,
        dateStr,
        value: val,
        dayOfWeek: dow,
      });
    }

    return {
      cells: list,
      totalMinutes: sum,
      activeDays: active,
    };
  }, [year, month, data]);

  const hoveredIntensity = hoveredDay ? getIntensityInfo(hoveredDay.item.value) : null;

  return (
    <div className="space-y-4">
      {/* Сводка за месяц */}
      <div className="flex flex-wrap items-center justify-between gap-4 text-sm text-zinc-400">
        <div className="flex items-center gap-6">
          <div>
            Активных дней в месяце:{' '}
            <span className="font-semibold text-zinc-100">{activeDays}</span>
          </div>
          <div>
            Всего за {MONTH_NAMES[month].toLowerCase()}:{' '}
            <span className="font-semibold text-emerald-400">
              {totalMinutes} мин ({(totalMinutes / 60).toFixed(1)} ч)
            </span>
          </div>
        </div>
      </div>

      {/* Заголовки дней недели */}
      <div className="grid grid-cols-7 gap-2 text-center text-xs font-medium text-zinc-500">
        {SHORT_WEEKDAYS.map((dayName) => (
          <div key={dayName} className="py-1">
            {dayName}
          </div>
        ))}
      </div>

      {/* Календарная сетка крупных ячеек */}
      <div className="grid grid-cols-7 gap-2">
        {cells.map((item, idx) => {
          if (!item) {
            return (
              <div
                key={`empty-${idx}`}
                className="h-16 rounded-xl bg-zinc-950/30 border border-zinc-900/50"
              />
            );
          }

          const intensity = getIntensityInfo(item.value);
          const isSelected = selectedDate === item.dateStr;
          const progressPct = Math.min(100, Math.round((item.value / 120) * 100));

          return (
            <button
              key={item.dateStr}
              type="button"
              onClick={() => onSelectDate?.(item.dateStr)}
              onMouseEnter={(e) =>
                setHoveredDay({ item, x: e.clientX, y: e.clientY })
              }
              onMouseMove={(e) =>
                setHoveredDay((prev) =>
                  prev ? { ...prev, x: e.clientX, y: e.clientY } : null
                )
              }
              onMouseLeave={() => setHoveredDay(null)}
              className={`group relative h-16 p-2.5 rounded-xl border text-left flex flex-col justify-between transition-all duration-150 hover:scale-[1.04] hover:z-10 ${
                intensity.cellClass
              } ${
                isSelected
                  ? 'ring-2 ring-emerald-400 ring-offset-2 ring-offset-zinc-950 z-10'
                  : ''
              }`}
            >
              <div className="flex items-center justify-between w-full">
                <span className="text-xs font-bold opacity-90">{item.dayNumber}</span>
                {item.value > 0 && (
                  <span className="text-[11px] font-mono opacity-95">
                    {item.value}м
                  </span>
                )}
              </div>

              {/* Мини прогресс-бар внутри дня */}
              <div className="w-full h-1.5 rounded-full bg-zinc-950/40 overflow-hidden">
                <div
                  style={{ width: `${progressPct}%` }}
                  className="h-full bg-current transition-all duration-300 rounded-full opacity-80"
                />
              </div>
            </button>
          );
        })}
      </div>

      {/* Всплывающий тултип */}
      {hoveredDay && hoveredIntensity && (
        <div
          style={{
            left: hoveredDay.x,
            top: hoveredDay.y - 14,
          }}
          className="fixed z-50 pointer-events-none -translate-x-1/2 -translate-y-full px-3.5 py-2.5 rounded-xl bg-zinc-950/95 border border-zinc-700/80 shadow-2xl backdrop-blur-md space-y-1.5 min-w-[200px]"
        >
          <div className="text-xs font-semibold text-zinc-100">
            {formatHumanDate(hoveredDay.item.dateStr)}
          </div>
          <div className="text-xs text-zinc-300 font-medium">
            {formatDuration(hoveredDay.item.value)}
          </div>
          <div>
            <span
              className={`inline-block text-[10px] font-medium px-2 py-0.5 rounded-full border ${hoveredIntensity.badgeClass}`}
            >
              {hoveredIntensity.label}
            </span>
          </div>
        </div>
      )}
    </div>
  );
};