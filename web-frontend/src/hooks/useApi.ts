import { useAuth } from '@clerk/clerk-react';
import axios, { type AxiosInstance } from 'axios';
import { create } from 'zustand';

type ApiClient = {
  instance: AxiosInstance;
};

const API_BASE_URL = import.meta.env.VITE_BASE_URL;

const useApiStore = create<ApiClient>(() => ({
  instance: axios.create({
    baseURL: API_BASE_URL,
    timeout: 10000,
    headers: {
      'Content-Type': 'application/json',
    },
  }),
}));

export const useApi = () => {
  const { getToken } = useAuth();
  const { instance } = useApiStore();

  // 設置請求攔截器
  instance.interceptors.request.use(
    async (config) => {
      try {
        const token = await getToken();
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      } catch (error) {
        console.error('Error getting token:', error);
        return config;
      }
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // 設置響應攔截器
  instance.interceptors.response.use(
    (response) => response,
    async (error) => {
      if (error.response?.status === 401) {
        // 處理認證錯誤
        console.error('Authentication error');
        // 這裡可以添加重新登入的邏輯
        // router.navigate({
        //   to: "/login",
        // });
      }
      return Promise.reject(error);
    }
  );

  return instance;
};
