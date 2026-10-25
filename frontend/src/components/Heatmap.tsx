import React, { useMemo, useState } from 'react';
import { ActivityGraph } from '../types/habit';
import {
  formatDuration,
  formatHumanDate,
  getIntensityInfo,
  SHORT_WEEKDAYS,
} from '../utils/date';

interface HeatmapProps {
  year: number;
  data: ActivityGraph;
  selectedDate?: string;
  onSelectDate?: (dateStr: string) => void;
}

interface DayCell {
  dateStr: string;
  value: number;
  dayOfWeek: number;
  weekIndex: number;
}

interface HoveredCellState {
  day: DayCell;
  clientX: number;
  clientY: number;
}

export const Heatmap: React.FC<HeatmapProps> = ({
  year,
  data,
  selectedDate,
  onSelectDate,
}) => {
  const [hovered, setHovered] = useState<HoveredCellState | null>(null);

  const { weeks, totalMinutes, activeDays } = useMemo(() => {
    const start = new Date(Date.UTC(year, 0, 1));
    const end = new Date(Date.UTC(year, 11, 31));

    const resultWeeks: (DayCell | null)[][] = [];
    let currentWeek: (DayCell | null)[] = Array(7).fill(null);
    let weekIdx = 0;

    const current = new Date(start);
    let isFirst = true;

    while (current <= end) {
      const dateStr = current.toISOString().split('T')[0];
      const dayOfWeek = (current.getUTCDay() + 6) % 7; // Пн = 0, Вс = 6

      if (!isFirst && dayOfWeek === 0) {
        resultWeeks.push(currentWeek);
        currentWeek = Array(7).fill(null);
        weekIdx++;
      }

      currentWeek[dayOfWeek] = {
        dateStr,
        value: data[dateStr] || 0,
        dayOfWeek,
        weekIndex: weekIdx,
      };

      isFirst = false;
      current.setUTCDate(current.getUTCDate() + 1);
    }
    resultWeeks.push(currentWeek);

    const values = Object.entries(data)
      .filter(([k]) => k.startsWith(`${year}-`))
      .map(([, v]) => v);

    const sumMinutes = values.reduce((acc, v) => acc + v, 0);
    const countActive = values.filter((v) => v > 0).length;

    return {
      weeks: resultWeeks,
      totalMinutes: sumMinutes,
      activeDays: countActive,
    };
  }, [year, data]);

  const hoveredIntensity = hovered ? getIntensityInfo(hovered.day.value) : null;

  return (
    <div className="space-y-4 relative">
      {/* Статистика сверху */}
      <div className="flex flex-wrap items-center justify-between gap-4 text-sm text-zinc-400">
        <div className="flex items-center gap-6">
          <div>
            Активных дней: <span className="font-semibold text-zinc-100">{activeDays}</span>
          </div>
          <div>
            Всего за {year} г.:{' '}
            <span className="font-semibold text-emerald-400">
              {totalMinutes} мин ({(totalMinutes / 60).toFixed(1)} ч)
            </span>
          </div>
        </div>

        <div className="flex items-center gap-1.5 text-xs">
          <span>Меньше</span>
          <span className="h-3 w-3 rounded-sm bg-zinc-900 border border-zinc-800" />
          <span className="h-3 w-3 rounded-sm bg-emerald-950 border border-emerald-800" />
          <span className="h-3 w-3 rounded-sm bg-emerald-800" />
          <span className="h-3 w-3 rounded-sm bg-emerald-600" />
          <span className="h-3 w-3 rounded-sm bg-emerald-400" />
          <span>Больше</span>
        </div>
      </div>

      {/* Сетка года с подписями дней недели */}
      <div className="overflow-x-auto pb-2">
        <div className="inline-flex gap-2 min-w-max items-start">
          {/* Ось Y: Дни недели */}
          <div className="flex flex-col gap-1 pr-1 pt-0.5 text-[10px] text-zinc-500 select-none">
            {SHORT_WEEKDAYS.map((dayName, idx) => (
              <div
                key={dayName}
                className={`h-3.5 flex items-center transition-colors ${
                  hovered?.day.dayOfWeek === idx ? 'text-emerald-400 font-semibold' : ''
                }`}
              >
                {idx % 2 === 0 ? dayName : ''}
              </div>
            ))}
          </div>

          {/* Колонки недель */}
          <div className="inline-flex gap-1">
            {weeks.map((week, wIdx) => {
              const isHoveredCol = hovered?.day.weekIndex === wIdx;
              return (
                <div
                  key={wIdx}
                  className={`flex flex-col gap-1 rounded-sm transition-colors ${
                    isHoveredCol ? 'bg-zinc-800/30' : ''
                  }`}
                >
                  {week.map((day, dIdx) => {
                    if (!day) {
                      return <div key={`empty-${wIdx}-${dIdx}`} className="h-3.5 w-3.5" />;
                    }
                    const intensity = getIntensityInfo(day.value);
                    const isSelected = selectedDate === day.dateStr;
                    const isHoveredRow = hovered?.day.dayOfWeek === dIdx;

                    return (
                      <button
                        key={day.dateStr}
                        type="button"
                        onClick={() => onSelectDate?.(day.dateStr)}
                        onMouseEnter={(e) =>
                          setHovered({
                            day,
                            clientX: e.clientX,
                            clientY: e.clientY,
                          })
                        }
                        onMouseMove={(e) =>
                          setHovered((prev) =>
                            prev
                              ? {
                                  ...prev,
                                  clientX: e.clientX,
                                  clientY: e.clientY,
                                }
                              : null
                          )
                        }
                        onMouseLeave={() => setHovered(null)}
                        className={`h-3.5 w-3.5 rounded-[3px] border transition-all duration-150 hover:scale-125 hover:z-20 hover:border-emerald-300 ${
                          intensity.cellClass
                        } ${isSelected ? 'ring-2 ring-emerald-400 ring-offset-1 ring-offset-zinc-950 z-10' : ''} ${
                          isHoveredRow && !hovered?.day.dateStr.includes(day.dateStr)
                            ? 'brightness-110'
                            : ''
                        }`}
                      />
                    );
                  })}
                </div>
              );
            })}
          </div>
        </div>
      </div>

      {/* Плавающий интерактивный тултип */}
      {hovered && hoveredIntensity && (
        <div
          style={{
            left: hovered.clientX,
            top: hovered.clientY - 14,
          }}
          className="fixed z-50 pointer-events-none -translate-x-1/2 -translate-y-full px-3.5 py-2.5 rounded-xl bg-zinc-950/95 border border-zinc-700/80 shadow-2xl backdrop-blur-md space-y-1.5 min-w-[200px] animate-in fade-in zoom-in-95 duration-100"
        >
          <div className="text-xs font-semibold text-zinc-100">
            {formatHumanDate(hovered.day.dateStr)}
          </div>
          <div className="text-xs text-zinc-300 font-medium">
            {formatDuration(hovered.day.value)}
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