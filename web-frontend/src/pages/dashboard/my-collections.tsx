import {
  FolderKanban,
  Share2,
  Clock,
  UserRound,
  ChartColumnBig,
  Eye,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import LanguageDistribution from '@/pages/dashboard/language-distribution';
import { UpdatedRepo } from '@/components/shared/repo';
import useCrontab from '@/apis/user/useCrontab';
import { formatDistance } from 'date-fns';
import { useQuery } from '@tanstack/react-query';
import { useApi } from '@/hooks/useApi';
import { useCallback } from 'react';

export default function MyCollections() {
  const api = useApi();

  const getUserLatestUpdated = useCallback(async () => {
    const { data } = await api.get<any[]>(`/repository/latest-updated`);
    return data;
  }, [api]);

  const { data, isLoading } = useQuery({
    queryKey: ['latest-updated'],
    queryFn: () => getUserLatestUpdated(),
  });

  return (
    <div className='row-span-5 grid grid-cols-3 h-full gap-4'>
      <div className='col-span-2 rounded-lg space-y-4 h-full'>
        <div className='flex gap-x-4 items-center'>
          <ChartColumnBig className='text-primary' />
          <div className='flex flex-col'>
            <h2 className='text-xl font-bold'>Overview</h2>
            <div className='text-slate-500 text-sm'>
              Your GitHub stars at a glance
            </div>
          </div>
        </div>
        <div className='grid grid-cols-2 grid-rows-1 gap-4'>
          <LanguageDistribution />
          <Block />
        </div>
      </div>
      <div className='col-span-1 rounded-lg space-y-4 h-full'>
        <div className='flex gap-x-4 items-center'>
          <Eye className='text-primary' />
          <div className='flex flex-col'>
            <h2 className='text-xl font-bold'>Recently Updated</h2>
            <div className='text-slate-500 text-sm'>
              See what’s new in your starred repos.
            </div>
          </div>
        </div>
        <div className='grid grid-cols-1 grid-rows-4 gap-4 h-[320px]'>
          {data &&
            data.length > 0 &&
            data
              .slice(0, 4)
              .map((repo) => <UpdatedRepo key={repo.id} {...repo} />)}
        </div>
      </div>
    </div>
  );
}

const Block = () => {
  const { data: crontab, isLoading } = useCrontab();

  return (
    <div className='grid grid-cols-2 grid-rows-3 gap-3'>
      <div className='flex flex-col gap-y-3 p-4 bg-slate-300 rounded-lg justify-center'>
        <div className='flex items-center gap-x-2 justify-between'>
          <div className='text-sm text-slate-900'>Last Synced</div>
          <Clock className='w-5 h-5' />
        </div>
        <div className='flex justify-start text-lg font-semibold'>
          {crontab &&
            formatDistance(crontab.updated_at, new Date(), { addSuffix: true })}
        </div>
      </div>
      <div className='flex flex-col gap-y-3 p-4 bg-slate-300 rounded-lg justify-center'>
        <div className='flex items-center gap-x-2 justify-between'>
          <div className='text-sm text-slate-900'>Friends</div>
          <UserRound className='w-5 h-5' />
        </div>
        <div className='flex justify-start text-lg font-semibold'>14</div>
      </div>
      <div className='flex flex-col gap-y-3 p-4 bg-slate-300 rounded-lg justify-center'>
        <div className='flex items-center gap-x-2 justify-between'>
          <div className='text-sm text-slate-900'>Your Collections</div>
          <FolderKanban className='w-5 h-5' />
        </div>
        <div className='flex justify-start text-lg font-semibold'>
          2 hour ago
        </div>
      </div>
      <div className='flex flex-col gap-y-3 p-4 bg-slate-300 rounded-lg justify-center'>
        <div className='flex items-center gap-x-2 justify-between'>
          <div className='text-sm text-slate-900'>Shared With You</div>
          <Share2 className='w-5 h-5' />
        </div>
        <div className='flex justify-start text-lg font-semibold'>
          2 hour ago
        </div>
      </div>
      <div className='flex items-end'>
        <Button variant='outline' className='w-full'>
          Trending Repository
        </Button>
      </div>
      <div className='flex items-end'>
        <Button variant='outline' className='w-full'>
          Trending Developer
        </Button>
      </div>
    </div>
  );
};
