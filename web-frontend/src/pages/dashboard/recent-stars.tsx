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
import GridRepo from '@/components/shared/repo';

export default function RecentStars() {
  return (
    <div className='row-span-3 grid grid-rows-1 h-full w-full'>
      <div className='flex flex-col gap-y-4 w-full'>
        <div className='flex gap-x-4 items-center'>
          <Star className='text-primary' />
          <div className='flex flex-col'>
            <h2 className='text-xl font-bold'>Latest Stars</h2>
            <div className='text-slate-500 text-sm'>
              Your recently starred GitHub repositories
            </div>
          </div>
        </div>
        <div className='grid grid-cols-3 gap-4'>
          <GridRepo />
          <GridRepo />
          <GridRepo />
        </div>
      </div>
    </div>
  );
}
