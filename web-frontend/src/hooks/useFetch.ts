import axios, { type AxiosRequestConfig, type HttpStatusCode } from 'axios';
import { useCallback, useEffect, useRef, useState } from 'react';
import { useApi } from './useApi';

type TError = {
  message: string;
  code: HttpStatusCode;
};

interface UseFetchProps<TData> {
  config: AxiosRequestConfig;
  initialRun?: boolean;
  onSuccess?: (data: TData) => void;
  onError?: (error: TError) => void;

  retry?: {
    count: number;
    delay: number;
  };

  cache?: {
    enabled: boolean;
    ttl: number;
  };
}

interface UseFetchState<TData> {
  data: TData | null;
  error: TError | null;
  isLoading: boolean;
  statusCode: number | null;
}

interface RunProps {
  params?: Record<string, any>;
  payload?: any;
}

class CacheManager {
  private cache: Map<string, { data: any; timestamp: number }> = new Map();

  set(key: string, data: any, ttl: number): void {
    this.cache.set(key, {
      data,
      timestamp: Date.now() + ttl,
    });
  }

  get(key: string): any | null {
    const item = this.cache.get(key);
    if (!item) return null;
    if (Date.now() > item.timestamp) {
      this.cache.delete(key);
      return null;
    }
    return item.data;
  }

  clear(): void {
    this.cache.clear();
  }
}

const cacheManager = new CacheManager();

export default function useFetch<TData>({
  config,
  initialRun = false,
  onSuccess,
  onError,
  retry = { count: 0, delay: 1000 },
  cache = { enabled: false, ttl: 5 * 60 * 1000 },
}: UseFetchProps<TData>) {
  const [state, setState] = useState<UseFetchState<TData>>({
    data: null,
    error: null,
    isLoading: initialRun,
    statusCode: null,
  });
  const api = useApi();

  const abortControllerRef = useRef<AbortController | null>(null);
  const retryCountRef = useRef(0);

  // 生成快取 key
  const getCacheKey = (params?: Record<string, any>, payload?: any): string => {
    return JSON.stringify({
      url: config.url,
      method: config.method,
      params,
      payload,
    });
  };

  const run = useCallback(
    async ({ params, payload }: RunProps = {}) => {
      // 取消之前的請求
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }

      abortControllerRef.current = new AbortController();

      // 檢查快取
      if (cache.enabled && config.method?.toLowerCase() === 'get') {
        const cacheKey = getCacheKey(params, payload);
        const cachedData = cacheManager.get(cacheKey);
        if (cachedData) {
          setState((prev) => ({
            ...prev,
            data: cachedData,
            isLoading: false,
          }));
          onSuccess?.(cachedData);
          return;
        }
      }

      setState((prev) => ({
        ...prev,
        isLoading: true,
        error: null,
        statusCode: null,
      }));

      const executeRequest = async (): Promise<void> => {
        try {
          const response = await api.request<TData>({
            ...config,
            url: config.url,
            method: config.method,
            params,
            data: payload,
            signal: abortControllerRef.current?.signal,
          });

          // 存入快取
          if (cache.enabled && config.method?.toLowerCase() === 'get') {
            const cacheKey = getCacheKey(params, payload);
            cacheManager.set(cacheKey, response.data, cache.ttl);
          }

          setState((prev) => ({
            ...prev,
            data: response.data,
            statusCode: response.status,
            isLoading: false,
          }));

          onSuccess?.(response.data);
          retryCountRef.current = 0;
        } catch (err) {
          if (axios.isAxiosError(err)) {
            if (retryCountRef.current < retry.count) {
              retryCountRef.current++;
              await new Promise((resolve) => setTimeout(resolve, retry.delay));
              return executeRequest();
            }

            const error = {
              message: err.response?.data?.error || err.message,
              code: err.response?.status,
            } as TError;

            setState((prev) => ({
              ...prev,
              error,
              statusCode: err.response?.status || null,
              isLoading: false,
            }));

            onError?.(error);
          }
        }
      };

      await executeRequest();
    },
    [config, onSuccess, onError, retry, cache]
  );

  useEffect(() => {
    if (initialRun) {
      run({
        params: config.params,
        payload: config.data,
      });
    }

    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, []);

  const reset = useCallback(() => {
    setState({
      data: null,
      error: null,
      isLoading: false,
      statusCode: null,
    });
  }, []);

  return {
    ...state,
    run,
    reset,
  };
}
