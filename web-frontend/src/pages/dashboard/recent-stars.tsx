import { Star } from 'lucide-react';
import { DetailRepo } from '@/components/shared/repo';
import { useQuery } from '@tanstack/react-query';
import { useApi } from '@/hooks/useApi';
import { useCallback } from 'react';

export default function RecentStars() {
  const api = useApi();

  const getUserLatestStarred = useCallback(async () => {
    const { data } = await api.get<any[]>(`/repository/latest-starred`);
    return data;
  }, [api]);

  const { data, isLoading } = useQuery({
    queryKey: ['latest-starred'],
    queryFn: () => getUserLatestStarred(),
  });

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
          {data &&
            data.length > 0 &&
            data
              .slice(0, 3)
              .map((repo) => <DetailRepo key={repo.id} {...repo} />)}
        </div>
      </div>
    </div>
  );
}
