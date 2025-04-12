import { motion } from 'framer-motion';
import clsx from 'clsx';
import { Button } from '@/components/ui/button';
import {
  Star,
  Eye,
  LucideCalendarRange,
  ExternalLink,
  DoorOpen,
  CircleDot,
  Users,
  GitFork,
} from 'lucide-react';
import { getLanguageColor } from '@/pages/dashboard/color';
import {
  daysAgo,
  Issues,
  Language,
  Stars,
  Forks,
} from '@/components/shared/repo/util';
import { cn } from '@/lib/utils';

function UpdatedRepo() {
  return (
    <motion.div
      initial={{ opacity: 0, y: -100 }}
      animate={{ opacity: 1, y: 0, transition: { delay: 1 * 0.05 } }}
      className={cn(
        'bg-white shadow-md dark:bg-slate-700/30 p-3 flex flex-col rounded-md gap-y-2'
      )}
    >
      <div className='flex gap-x-3 items-center'>
        <img
          src='https://avatars.githubusercontent.com/u/50438175?v=4'
          alt='encore'
          width={20}
          height={20}
          className='rounded-full'
        />
        <div>
          <span>encoredev</span>
          <span>/</span>
          <span className='font-bold'>encore</span>
        </div>
      </div>
      <div className='flex w-full justify-between'>
        <Language />
        <div className='text-sm text-slate-400'>Update 2 hours ago</div>
      </div>
    </motion.div>
  );
}

export { UpdatedRepo };
