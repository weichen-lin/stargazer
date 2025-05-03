import { useApi } from '@/hooks/useApi';
import { useQuery } from '@tanstack/react-query';
import { useCallback } from 'react';

interface CrontabDto {
  stargazers: number;
  created_at: string;
  updated_at: string;
}

export default function useCrontab() {
  const api = useApi();

  const getUserCrontab = useCallback(async () => {
    const { data } = await api.get<CrontabDto>(`/user/crontab`);
    return data;
  }, []);

  const { data, isLoading } = useQuery({
    queryKey: ['pings'],
    queryFn: () => getUserCrontab(),
  });

  return { data, isLoading };
}
