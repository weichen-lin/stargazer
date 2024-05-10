'use client'

import { create } from 'zustand'
import { useState, useEffect } from 'react'
import { IUserConfig, getUserInfo, updateInfo } from '@/actions/neo4j'
import { useUser } from '@/context'

interface IUseConfig {
  config: IUserConfig
  setConfig: (config: IUserConfig) => void
  change: (key: keyof IUserConfig, value: string | number) => void
}

export const useConfigState = create<IUseConfig>(set => ({
  config: {
    openAIKey: '',
    githubToken: '',
    limit: 5,
    cosine: 0.8,
  },
  setConfig: (config: IUserConfig) => set(state => ({ config: { ...state.config, ...config } })),
  change: (key, value) => set(state => ({ config: { ...state.config, [key]: value } })),
}))

const useConfig = () => {
  const { config, setConfig, change } = useConfigState()
  const [isLoading, setIsLoading] = useState(true)
  const { email } = useUser()

  const update = async () => {
    try {
      setIsLoading(true)
      await updateInfo({ email, ...config })
      setConfig({ ...config })
    } catch {
      console.error('failed to update user setting')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    const fetchSetting = async () => {
      setIsLoading(true)
      const info = await getUserInfo(email)
      setConfig(info)
      setIsLoading(false)
    }

    fetchSetting()
  }, [])

  return { isLoading, config, change, update }
}

export default useConfig
