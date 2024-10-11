'use client'

import { useState, useEffect, useCallback } from 'react'
import axios, { AxiosRequestConfig, AxiosError } from 'axios'

interface useFetchProps<T> {
  config: AxiosRequestConfig<T>
  initialRun?: boolean
  onSuccess?: (res: T) => void
  onError?: (err: any) => void
}

const instance = axios.create({
  baseURL: '/api',
  validateStatus: statusCode => {
    return statusCode < 300
  },
})

export interface IRunProps {
  params?: Record<string, string>
  payload?: any
}

export default function useFetch<T>(props: useFetchProps<T>) {
  const { config, initialRun, onSuccess, onError } = props
  const { url, method, params, data: payload, headers } = config

  const [isLoading, setIsLoading] = useState(initialRun ? true : false)
  const [error, setError] = useState<any>(null)
  const [statusCode, setStatusCode] = useState<number | null>(null)
  const [data, setData] = useState<T | null>(null)

  const run = useCallback(async ({ params, payload }: IRunProps) => {
    try {
      setIsLoading(true)
      setError(null)
      setStatusCode(null)

      const response = await instance.request<T>({
        url,
        method,
        params,
        data: payload,
        headers,
      })

      setData(response.data)
      setStatusCode(response.status)

      onSuccess && onSuccess(response.data)
    } catch (err) {
      if (err instanceof AxiosError) {
        onError && onError(err.response?.data)
      }
    } finally {
      setIsLoading(false)
    }
  }, [])

  useEffect(() => {
    if (initialRun) run({ params, payload })
  }, [])

  return { statusCode, error, isLoading, data, run }
}
