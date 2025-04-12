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

export default function GridRepo() {
  return (
    <motion.div
      initial={{ opacity: 0, y: -100 }}
      animate={{ opacity: 1, y: 0, transition: { delay: 1 * 0.05 } }}
      className={clsx(
        'shadow-md bg-white dark:bg-slate-700/30 p-3 flex flex-col rounded-md h-full',
        'border-[1px] dark:border-slate-100/30 border-slate-500/10 relative md:overflow-hidden'
      )}
      // onHoverStart={() => setIsHover(true)}
      // onHoverEnd={() => setIsHover(false)}
    >
      <div className='flex flex-col justify-between gap-y-2 h-full'>
        <div className='flex gap-x-3 items-center'>
          <img
            src='https://avatars.githubusercontent.com/u/50438175?v=4'
            alt='encore'
            width={40}
            height={40}
            className='rounded-full'
          />
          <div>
            <span>encoredev</span>
            <span>/</span>
            <span className='font-bold'>encore</span>
          </div>
        </div>
        <div className='line-clamp-2 text-sm text-slate-500/75 dark:text-white/70'>
          Open Source Development Platform for building robust type-safe
          distributed systems with declarative infrastructure
        </div>
        <div className='flex justify-between items-center mb-4'>
          <div className='flex gap-x-2 md:hidden'>
            <a href='https://github.com/encoredev/encore' target='_blank'>
              <Button size='sm' variant='outline'>
                <ExternalLink className='w-4 h-4' />
              </Button>
            </a>
            <Button size='sm' variant='outline'>
              <DoorOpen className='w-4 h-4' />
            </Button>
          </div>
          {/* <DayAgo dateString={updated_at} /> */}
        </div>
        <div className='grid grid-cols-5'>
          <Language />
          <Issues count={10} />
          <Stars count={100} />
          <Forks count={50} />
        </div>
      </div>
      {/* {isHover && (
        <motion.div
          initial={{ opacity: 0, y: 170, scale: 0.9 }}
          animate={{
            opacity: 1,
            y: 90,
            scale: 1,
            transition: { duration: 0.2 },
          }}
          className={clsx(
            'absolute inset-0 bg-opacity-20 flex items-center justify-center gap-x-4 w-full h-20 backdrop-blur-sm',
            'from-slate-200/40 to-slate-100/10 dark:from-slate-700 dark:to-slate-800 bg-gradient-to-t'
          )}
        >
          <a href={html_url} target='_blank'>
            <Button size='sm' variant='outline'>
              <ExternalLink className='w-4 h-4' />
            </Button>
          </a>
          <Button
            size='sm'
            variant='outline'
            onClick={() => {
              setRepoID(repo_id);
              setOpen(true);
            }}
          >
            <DoorOpen className='w-4 h-4' />
          </Button>
        </motion.div>
      )} */}
    </motion.div>
  );
}
