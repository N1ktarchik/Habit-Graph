export const MONTH_NAMES = [
  'Январь',
  'Февраль',
  'Март',
  'Апрель',
  'Май',
  'Июнь',
  'Июль',
  'Август',
  'Сентябрь',
  'Октябрь',
  'Ноябрь',
  'Декабрь',
];

export const MONTH_NAMES_GENITIVE = [
  'января',
  'февраля',
  'марта',
  'апреля',
  'мая',
  'июня',
  'июля',
  'августа',
  'сентября',
  'октября',
  'ноября',
  'декабря',
];

export const SHORT_WEEKDAYS = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'];

export function formatHumanDate(dateStr: string): string {
  const [y, m, d] = dateStr.split('-').map(Number);
  if (!y || !m || !d) return dateStr;
  const date = new Date(Date.UTC(y, m - 1, d));
  const dayOfWeek = (date.getUTCDay() + 6) % 7;
  return `${d} ${MONTH_NAMES_GENITIVE[m - 1]} ${y}, ${SHORT_WEEKDAYS[dayOfWeek]}`;
}

export function formatDuration(minutes: number): string {
  if (minutes <= 0) return 'Нет активности';
  const hrs = Math.floor(minutes / 60);
  const mins = minutes % 60;
  if (hrs === 0) return `${mins} мин`;
  if (mins === 0) return `${minutes} мин (${hrs} ч)`;
  return `${minutes} мин (${hrs} ч ${mins} мин)`;
}

export function formatTimerSeconds(totalSeconds: number): string {
  const hrs = Math.floor(totalSeconds / 3600);
  const mins = Math.floor((totalSeconds % 3600) / 60);
  const secs = totalSeconds % 60;
  return [hrs, mins, secs].map((v) => String(v).padStart(2, '0')).join(':');
}

export interface IntensityInfo {
  label: string;
  badgeClass: string;
  cellClass: string;
  barClass: string;
}

export function getIntensityInfo(minutes: number): IntensityInfo {
  if (minutes <= 0) {
    return {
      label: 'Нет активности',
      badgeClass: 'bg-zinc-800 text-zinc-400 border-zinc-700',
      cellClass: 'bg-zinc-900/90 border-zinc-800/80 hover:border-zinc-600',
      barClass: 'bg-zinc-800',
    };
  }
  if (minutes < 30) {
    return {
      label: 'Разминка (< 30м)',
      badgeClass: 'bg-emerald-950/90 text-emerald-300 border-emerald-800',
      cellClass: 'bg-emerald-950 border-emerald-800/70 text-emerald-300',
      barClass: 'bg-emerald-800',
    };
  }
  if (minutes < 60) {
    return {
      label: 'Хороший темп (30–60м)',
      badgeClass: 'bg-emerald-900/80 text-emerald-200 border-emerald-700',
      cellClass: 'bg-emerald-800/80 border-emerald-600/60 text-emerald-100',
      barClass: 'bg-emerald-600',
    };
  }
  if (minutes < 120) {
    return {
      label: 'Продуктивно (1–2ч)',
      badgeClass: 'bg-emerald-600/30 text-emerald-300 border-emerald-500/50',
      cellClass: 'bg-emerald-600 border-emerald-400/70 text-zinc-950 font-semibold',
      barClass: 'bg-emerald-500',
    };
  }
  return {
    label: 'Максимум фокуса (2ч+)',
    badgeClass: 'bg-emerald-400/20 text-emerald-300 border-emerald-400',
    cellClass:
      'bg-emerald-400 border-emerald-300 text-zinc-950 font-bold shadow-[0_0_10px_rgba(52,211,153,0.4)]',
    barClass: 'bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.5)]',
  };
}