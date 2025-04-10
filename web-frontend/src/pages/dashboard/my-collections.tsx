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

export default function MyCollections() {
  return (
    <div className='row-span-5 grid grid-cols-3 h-full gap-4'>
      <div className='col-span-2 rounded-lg space-y-4'>
        <div className='flex justify-between py-1'>
          <div className='flex gap-x-4 items-center'>
            <FolderKanban className='text-primary' />
            <div className='flex flex-col'>
              <h2 className='text-xl font-bold'>My Collections</h2>
              <div className='text-slate-500 text-sm'>
                Organize your starred repositories
              </div>
            </div>
          </div>
        </div>
        <div className='grid grid-rows-2 grid-flow-col gap-4'>
          <Card />
          <Card />
          <Card />
          <Card />
        </div>
      </div>
      <div className='col-span-1 rounded-lg'>
        <div className='flex justify-between py-1'>
          <div className='flex gap-x-4 items-center'>
            <FolderKanban className='text-primary' />
            <div className='flex flex-col'>
              <h2 className='text-xl font-bold'>Shared With Me</h2>
              <div className='text-slate-500 text-sm'>
                Collections shared by others
              </div>
            </div>
          </div>
        </div>
        <div className='grid grid-rows-2 grid-flow-col gap-4 mt-4'>
          <Card2 />
          <Card2 />
        </div>
      </div>
    </div>
  );
}

const Card = () => (
  <MagicCard
    gradientColor='#D9D9D955'
    className='p-4 bg-card shadow-lg rounded-lg'
  >
    <div className='flex flex-col gap-y-2 w-full h-full'>
      <div className='flex justify-between items-center w-full'>
        <div className='flex gap-x-3 items-center'>
          <div className='bg-background w-7 h-7 flex items-center rounded-md justify-center'>
            <LucideFolder className='w-4 h-4 text-primary' />
          </div>
          <h1 className='font-semibold'>Frontend Framework</h1>
          {true ? (
            <LockOpen className='w-4 h-4 text-green-700' />
          ) : (
            <LucideLock className='w-5 h-5 mt-[2px]' />
          )}
        </div>
        <div className='flex gap-x-2 items-center'>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant='ghost' size='icon' className=''>
                <div className='flex items-center gap-x-2'>
                  <EllipsisVertical className='w-12 h-12' />
                </div>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align='end'>
              <DropdownMenuItem
                onClick={() => {}}
                className='flex gap-x-3 items-center'
              >
                <PencilIcon className='w-4 h-4' />
                Rename
              </DropdownMenuItem>
              <DropdownMenuItem className='flex gap-x-3 items-center'>
                <Trash2Icon className='w-4 h-4' />
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
          <Button variant='ghost'>
            <div className='flex items-center gap-x-2'>
              <ArrowBigRight size={24} />
            </div>
          </Button>
        </div>
      </div>
      <div className='line-clamp-2 text-sm text-slate-600 min-h-[2.5rem]'>
        line-clamp-2line-clamp-2line-clamp-2line-clamp-2li ne-
      </div>
      <div className='flex justify-between'>
        <div className='text-slate-400 text-sm'>12 repositories</div>
        <div className='text-slate-400 text-sm'>update 2 hours ago</div>
      </div>
    </div>
  </MagicCard>
);

const Card2 = () => (
  <MagicCard
    gradientColor='#D9D9D955'
    className='p-4 bg-card shadow-lg rounded-lg'
  >
    <div className='flex flex-col gap-y-2 w-full h-full'>
      <div className='flex justify-between items-center w-full min-h-[36px]'>
        <div className='flex gap-x-3 items-center'>
          <div className='bg-background w-7 h-7 flex items-center rounded-md justify-center'>
            <LucideFolder className='w-4 h-4 text-primary' />
          </div>
          <h1 className='font-semibold'>Frontend Framework</h1>
          {true ? (
            <LockOpen className='w-4 h-4 text-green-700' />
          ) : (
            <LucideLock className='w-5 h-5 mt-[2px]' />
          )}
        </div>
        <div className='flex items-center gap-x-3'>
          <AnimatedTooltip />
          <Button variant='ghost' className='h-8 w-8 p-0'>
            <ArrowBigRight className='w-4 h-4' />
          </Button>
        </div>
      </div>
      <div className='line-clamp-2 text-sm text-slate-600 min-h-[2.5rem]'>
        ne-clamp-2line-clamp-2line-clamp-2line-clamp-2line-clamp-2line-clamp-2line-clamp-2
        line-clamp-2line-clamp-2line -clamp-2line-clamp-2lin
        e-clamp-2line-clamp-2lin e-clamp-2
      </div>
      <div className='flex justify-between'>
        <div className='text-slate-400 text-sm'>12 repositories</div>
      </div>
    </div>
  </MagicCard>
);
