import { useQuery } from '@tanstack/react-query';
import { useApi } from '@/hooks/useApi';
import { useCallback } from 'react';

export interface LanguageDistribution {
  language: string;
  count: number;
}

export default function useLanguageDistribution() {
  const api = useApi();

  const getLanguageDistribution = useCallback(async () => {
    const { data } = await api.get<LanguageDistribution[]>(
      `/repository/language-distribution`
    );
    return data;
  }, [api]);

  const { data, isLoading } = useQuery({
    queryKey: ['language-distribution'],
    queryFn: () => getLanguageDistribution(),
  });

  return { data: data ?? [], isLoading };
}
