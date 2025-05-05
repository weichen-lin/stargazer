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
import { getLanguageColor, type Language } from '@/pages/dashboard/color';
import { cn } from '@/lib/utils';

import type { Repository } from '@/components/shared/repo/type';

function UpdatedRepo(props: Repository) {
  const {
    id,
    name,
    owner_name,
    avatar_url,
    html_url,
    homepage,
    description,
    created_at,
    updated_at,
    synced_at,
    watchers,
    open_issues,
    forks,
    language,
    archived,
    topics,
  } = props;

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
          src={avatar_url}
          alt={`${owner_name} avatar`}
          width={20}
          height={20}
          className='rounded-full'
        />
        <div>
          <span>{owner_name}</span>
          <span>/</span>
          <span className='font-bold'>{name}</span>
        </div>
      </div>
      <div className='flex w-full justify-between'>
        <div className='col-span-2 gap-x-2 justify-start items-center flex'>
          <div
            className='rounded-full w-2 h-2'
            style={{
              backgroundColor: getLanguageColor(language as Language),
            }}
          ></div>
          <div className='text-slate-500/75 dark:text-white/70 text-sm'>
            {language}
          </div>
        </div>
        <div className='text-sm text-slate-400'>Update 2 hours ago</div>
      </div>
    </motion.div>
  );
}

export { UpdatedRepo };
