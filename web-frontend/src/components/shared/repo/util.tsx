import { Star, LucideCalendarRange, CircleDot, GitFork } from 'lucide-react';
import { getLanguageColor } from '@/pages/dashboard/color';

const Issues = ({ count }: { count: number }) => {
  return (
    <div className='col-span-1 gap-x-2 justify-end items-center flex'>
      <CircleDot className='h-4 w-4' />
      <span className='text-slate-700 dark:text-white/70 font-light text-sm'>
        {count}
      </span>
    </div>
  );
};

const Stars = ({ count }: { count: number }) => {
  return (
    <div className='col-span-1 gap-x-2 justify-end items-center flex'>
      <Star className='h-4 w-4' />
      <span className='text-slate-700 dark:text-white/70 font-light text-sm'>
        {count}
      </span>
    </div>
  );
};

const Forks = ({ count }: { count: number }) => {
  return (
    <div className='col-span-1 gap-x-2 justify-end items-center flex'>
      <GitFork className='h-4 w-4' />
      <span className='text-slate-700 dark:text-white/70 font-light text-sm'>
        {count}
      </span>
    </div>
  );
};

function daysAgo(dateString: string) {
  const date = new Date(dateString);
  const currentDate = new Date();
  const timeDiff = currentDate.getTime() - date.getTime();
  const daysDiff = Math.floor(timeDiff / (1000 * 3600 * 24));
  return daysDiff;
}

const DayAgo = ({ dateString }: { dateString: string }) => {
  const days = daysAgo(dateString);
  return (
    <div className='flex gap-x-2 items-center justify-end'>
      <LucideCalendarRange className='h-5 w-5 text-gray-500' />
      <span className='dark:text-white/70 text-sm text-gray-500'>
        {days} days ago
      </span>
    </div>
  );
};

const Language = () => {
  return (
    <div className='col-span-2 gap-x-2 justify-start items-center flex'>
      <div
        className='rounded-full w-2 h-2'
        style={{
          backgroundColor: getLanguageColor('Go'),
        }}
      ></div>
      <div className='text-slate-500/75 dark:text-white/70 text-sm'>Go</div>
    </div>
  );
};

export { Issues, Stars, Forks, DayAgo, Language, getLanguageColor, daysAgo };
