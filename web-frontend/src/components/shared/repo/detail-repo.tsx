import { motion } from 'framer-motion';
import clsx from 'clsx';
import { Button } from '@/components/ui/button';
import { ExternalLink, DoorOpen } from 'lucide-react';
import { Issues, Language, Stars, Forks } from '@/components/shared/repo/util';
import type { Repository } from '@/components/shared/repo/type';

export default function DetailRepo(props: Repository) {
  const {
    name,
    owner_name,
    avatar_url,
    description,
    watchers,
    open_issues,
    forks,
  } = props;

  return (
    <motion.div
      initial={{ opacity: 0, y: -100 }}
      animate={{ opacity: 1, y: 0, transition: { delay: 1 * 0.05 } }}
      className={clsx(
        'shadow-md bg-white dark:bg-slate-700/30 p-3 flex flex-col rounded-xl h-full',
        'border-[1px] dark:border-slate-100/30 border-slate-500/10 relative md:overflow-hidden'
      )}
    >
      <div className='flex flex-col justify-between gap-y-2 h-full'>
        <div className='flex gap-x-3 items-center'>
          <img
            src={avatar_url}
            alt={`${owner_name} avatar`}
            width={40}
            height={40}
            className='rounded-full'
          />
          <div>
            <span>{owner_name}</span>
            <span>/</span>
            <span className='font-bold'>{name}</span>
          </div>
        </div>
        <div className='line-clamp-2 text-sm text-slate-500/75 dark:text-white/70'>
          {description}
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
        </div>
        <div className='grid grid-cols-5'>
          <Language />
          <Issues count={open_issues} />
          <Stars count={watchers} />
          <Forks count={forks} />
        </div>
      </div>
    </motion.div>
  );
}
