import { createFileRoute } from '@tanstack/react-router';
import { CreateCollection } from '@/pages/collections';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  FolderKanban,
  Share2,
  LucideFolder,
  LockOpen,
  LucideLock,
  Trash2Icon,
  EllipsisVertical,
  PencilIcon,
  ArrowBigRight,
  Calendar,
  Database,
  Eye,
  EyeOff,
  Users,
  Edit,
  Trash2,
  UsersRound,
} from 'lucide-react';

import { MagicCard } from '@/components/shared/magic-card';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { AnimatedTooltip } from '@/components/ui/animated-tooltip';

export const Route = createFileRoute('/_auth/collections')({
  component: Collections,
});

function Collections() {
  return (
    <div className='grid grid-rows-[50px_1fr] h-full gap-6'>
      <div className='flex items-center justify-between py-4 gap-x-4'>
        <CreateCollection />
        <Input className='border-slate-700' />
        <Select>
          <SelectTrigger className='w-[180px] border-slate-700'>
            <SelectValue placeholder='Theme' />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value='light'>Light</SelectItem>
            <SelectItem value='dark'>Dark</SelectItem>
            <SelectItem value='system'>System</SelectItem>
          </SelectContent>
        </Select>
        <Select>
          <SelectTrigger className='w-[180px] border-slate-700'>
            <SelectValue placeholder='Theme' />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value='light'>Light</SelectItem>
            <SelectItem value='dark'>Dark</SelectItem>
            <SelectItem value='system'>System</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className='grid grid-rows-3 grid-cols-3 grid-flow-col gap-4'>
        <Card />
        <Card />
        <Card />
        <Card />
        <Card />
        <Card />
        <Card />
        <Card />
      </div>
    </div>
  );
}

const Card = () => (
  <MagicCard
    gradientColor='#D9D9D955'
    className='p-4 bg-card shadow-lg rounded-lg'
  >
    <div className='flex flex-col gap-y-2 w-full h-full justify-between'>
      <div className='flex justify-between items-center w-full'>
        <h1 className='font-semibold'>Frontend Framework</h1>
        {true ? (
          <UsersRound className='w-4 h-4 text-green-700' />
        ) : (
          <LucideLock className='w-5 h-5 mt-[2px]' />
        )}
        <div className='flex gap-x-2 items-center'>
          <Button variant='ghost' className=''>
            <ArrowBigRight className='w-5 h-5' />
          </Button>
        </div>
      </div>
      <div className='line-clamp-2 text-sm text-slate-600 min-h-[2.5rem]'>
        line-clamp-2line-clamp-2line-clamp-2line-clamp-2li ne-
      </div>
      <div className='space-y-2'>
        <div className='flex gap-x-2 items-center'>
          <Database className='w-4 h-4 text-slate-700' />
          <div className='text-slate-700 text-sm'>12 repositories</div>
        </div>
        <div className='flex gap-x-2 items-center'>
          <Calendar className='w-4 h-4 text-slate-700' />
          <div className='text-slate-700 text-sm'>Update at 2 hours ago</div>
        </div>
        <div className='flex gap-x-2 items-center'>
          <UsersRound className='w-4 h-4 text-slate-700' />
          <div className='text-slate-700 text-sm'>Update at 2 hours ago</div>
        </div>
      </div>
    </div>
  </MagicCard>
);
