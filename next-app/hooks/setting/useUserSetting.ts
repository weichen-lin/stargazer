import { create } from 'zustand'
import { useState, useEffect } from 'react'
import { IUserSetting, getUserInfo, updateInfo } from '@/actions/neo4j'
import { useSession } from 'next-auth/react'

interface Setting extends IUserSetting {
  change: (key: keyof IUserSetting, value: string | number) => void
}

export const useSetting = create<Setting>(set => ({
  openAIKey: '',
  githubToken: '',
  limit: 5,
  cosine: 0.8,
  change: (key, value) => set(state => ({ ...state, [key]: value })),
}))

const useUserSetting = () => {
  const { openAIKey, githubToken, limit, cosine, change } = useSetting()
  const [isLoading, setIsLoading] = useState(true)
  const session = useSession()
  const email = session.data?.user?.email ?? ''

  const update = async () => {
    try {
      setIsLoading(true)
      await updateInfo({ email, githubToken, limit, openAIKey, cosine })
      change('openAIKey', openAIKey)
      change('limit', limit)
      change('cosine', cosine)
      change('githubToken', githubToken)
    } catch {
      console.error('failed to update user setting')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    const fetchSetting = async () => {
      setIsLoading(true)
      const info = await getUserInfo({ email })
      if (info) {
        change('openAIKey', info.openAIKey ?? '')
        change('limit', info.limit ?? 5)
        change('cosine', info?.cosine ?? 0.8)
        change('githubToken', info.githubToken ?? '')
      }
      setIsLoading(false)
    }
    if (email) {
      fetchSetting()
    }
  }, [])

  return { isLoading, openAIKey, githubToken, limit, cosine, change, update }
}

export default useUserSetting
