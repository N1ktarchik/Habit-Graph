import React, { useMemo, useState } from 'react';
import { TrendingUp, Award } from 'lucide-react';
import { ActivityGraph } from '../types/habit';
import {
  formatDuration,
  formatHumanDate,
  getIntensityInfo,
  MONTH_NAMES,
} from '../utils/date';

interface MonthTrendChartProps {
  year: number;
  month: number; // 0 .. 11
  data: ActivityGraph;
  selectedDate?: string;
  onSelectDate?: (dateStr: string) => void;
}

interface ChartPoint {
  day: number;
  dateStr: string;
  minutes: number;
}

export const MonthTrendChart: React.FC<MonthTrendChartProps> = ({
  year,
  month,
  data,
  selectedDate,
  onSelectDate,
}) => {
  const [hoveredIdx, setHoveredIdx] = useState<number | null>(null);

  const { points, maxMinutes, avgMinutes, peakPoint } = useMemo(() => {
    const daysInMonth = new Date(Date.UTC(year, month + 1, 0)).getUTCDate();
    const list: ChartPoint[] = [];
    let sum = 0;
    let peak: ChartPoint | null = null;

    for (let d = 1; d <= daysInMonth; d++) {
      const dateStr = new Date(Date.UTC(year, month, d)).toISOString().split('T')[0];
      const minutes = data[dateStr] || 0;
      const pt = { day: d, dateStr, minutes };
      list.push(pt);
      sum += minutes;
      if (!peak || minutes > peak.minutes) {
        peak = pt;
      }
    }

    const rawMax = Math.max(...list.map((p) => p.minutes), 60);
    // Округляем верхнюю границу графика до удобного числа (кратно 30)
    const niceMax = Math.ceil(rawMax / 30) * 30;
    const avg = daysInMonth > 0 ? Math.round(sum / daysInMonth) : 0;

    return {
      points: list,
      maxMinutes: niceMax,
      avgMinutes: avg,
      peakPoint: peak && peak.minutes > 0 ? peak : null,
    };
  }, [year, month, data]);

  const yTicks = [maxMinutes, Math.round(maxMinutes / 2), 0];

  // Генерация координат для плавной SVG-линии поверх столбцов
  const svgPolyline = useMemo(() => {
    if (points.length === 0) return '';
    return points
      .map((pt, idx) => {
        const x = ((idx + 0.5) / points.length) * 100;
        const y = 100 - (pt.minutes / maxMinutes) * 100;
        return `${x.toFixed(2)},${y.toFixed(2)}`;
      })
      .join(' ');
  }, [points, maxMinutes]);

  const hoveredPoint = hoveredIdx !== null ? points[hoveredIdx] : null;

  return (
    <div className="p-6 rounded-xl bg-zinc-900/60 border border-zinc-800/80 space-y-5">
      {/* Шапка графика со статистикой и пиками */}
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div className="space-y-0.5">
          <h3 className="text-sm font-semibold text-zinc-100 flex items-center gap-2">
            <TrendingUp className="h-4 w-4 text-emerald-400" />
            Динамика за {MONTH_NAMES[month].toLowerCase()} {year}
          </h3>
          <p className="text-xs text-zinc-400">
            Распределение времени по дням месяца и линия среднего темпа
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-4 text-xs">
          <div className="px-3 py-1.5 rounded-lg bg-zinc-950 border border-zinc-800 text-zinc-300">
            Среднее в день:{' '}
            <span className="font-semibold text-emerald-400">{avgMinutes} мин</span>
          </div>
          {peakPoint && (
            <div className="px-3 py-1.5 rounded-lg bg-emerald-500/10 border border-emerald-500/30 text-emerald-300 flex items-center gap-1.5">
              <Award className="h-3.5 w-3.5" />
              <span>
                Пик: {peakPoint.day}-е число ({peakPoint.minutes} мин)
              </span>
            </div>
          )}
        </div>
      </div>

      {/* Информационная плашка при наведении на столбец */}
      <div className="h-6 flex items-center justify-between text-xs px-2 rounded bg-zinc-950/60 border border-zinc-800/50">
        {hoveredPoint ? (
          <>
            <span className="text-zinc-300 font-medium">
              {formatHumanDate(hoveredPoint.dateStr)}
            </span>
            <span className="text-emerald-400 font-semibold">
              {formatDuration(hoveredPoint.minutes)}
            </span>
          </>
        ) : (
          <span className="text-zinc-500">
            Наведи курсор на любой день графика для просмотра точного времени
          </span>
        )}
      </div>

      {/* Область графика */}
      <div className="flex gap-3 items-stretch h-44 pt-2">
        {/* Ось Y */}
        <div className="flex flex-col justify-between text-[10px] text-zinc-500 font-mono pr-1 select-none text-right w-12">
          {yTicks.map((t) => (
            <span key={t}>{t}м</span>
          ))}
        </div>

        {/* Полотно графика */}
        <div className="relative flex-1 flex items-end gap-1 border-b border-l border-zinc-800 pl-1 pb-1">
          {/* Горизонтальная пунктирная линия среднего значения */}
          {avgMinutes > 0 && (
            <div
              style={{
                bottom: `${Math.min(100, (avgMinutes / maxMinutes) * 100)}%`,
              }}
              className="absolute left-0 right-0 border-t border-dashed border-emerald-500/40 pointer-events-none z-10"
              title={`Среднее: ${avgMinutes} мин`}
            />
          )}

          {/* SVG линия тренда */}
          <svg
            viewBox="0 0 100 100"
            preserveAspectRatio="none"
            className="absolute inset-0 w-full h-full pointer-events-none overflow-visible z-10"
          >
            <polyline
              fill="none"
              stroke="rgba(52, 211, 153, 0.55)"
              strokeWidth="1.5"
               vectorEffect="non-scaling-stroke"
              points={svgPolyline}
            />
          </svg>

          {/* Столбцы по дням */}
          {points.map((pt, idx) => {
            const heightPct = Math.max(
              pt.minutes > 0 ? 6 : 2,
              Math.round((pt.minutes / maxMinutes) * 100)
            );
            const intensity = getIntensityInfo(pt.minutes);
            const isPeak = peakPoint?.dateStr === pt.dateStr;
            const isSelected = selectedDate === pt.dateStr;

            return (
              <button
                key={pt.dateStr}
                type="button"
                onClick={() => onSelectDate?.(pt.dateStr)}
                onMouseEnter={() => setHoveredIdx(idx)}
                onMouseLeave={() => setHoveredIdx(null)}
                className="group relative flex-1 h-full flex flex-col justify-end items-center focus:outline-none z-20"
              >
                <div
                  style={{ height: `${heightPct}%` }}
                  className={`w-full rounded-t-sm transition-all duration-200 group-hover:brightness-125 ${
                    intensity.barClass
                  } ${isPeak ? 'ring-1 ring-emerald-300' : ''} ${
                    isSelected ? 'ring-2 ring-white' : ''
                  }`}
                />
              </button>
            );
          })}
        </div>
      </div>

      {/* Ось X: числа месяца */}
      <div className="flex gap-1 pl-14 text-[10px] text-zinc-500 font-mono select-none">
        {points.map((pt, idx) => {
          const showLabel =
            pt.day === 1 ||
            pt.day % 5 === 0 ||
            pt.day === points.length ||
            hoveredIdx === idx;
          return (
            <div
              key={pt.dateStr}
              className={`flex-1 text-center transition-colors ${
                hoveredIdx === idx ? 'text-emerald-400 font-bold' : ''
              }`}
            >
              {showLabel ? pt.day : ''}
            </div>
          );
        })}
      </div>
    </div>
  );
};