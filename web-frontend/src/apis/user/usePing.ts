import { useApi } from '@/hooks/useApi';
import { useQuery } from '@tanstack/react-query';
import { useCallback } from 'react';

export default function usePings() {
  const api = useApi();

  const getDebtTransactions = useCallback(async () => {
    const { data } = await api.get<any[]>(`/crontab/ping`);
    return data;
  }, []);

  const { data, isPending, isLoading } = useQuery({
    queryKey: ['pings'],
    queryFn: () => getDebtTransactions(),
  });

  const transactions = isPending || !data ? [] : data;
  const isLoding = isPending || isLoading;

  return { transactions, isLoding };
}
