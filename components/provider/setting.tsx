'use client'

import { UserInfo } from '@/actions/neo4j'
import { createContext, useState, useContext } from 'react'

interface Setting extends Pick<UserInfo, 'openAIKey' | 'githubToken' | 'limit' | 'cosine'> {
  changeKey: (key: string) => void
  changeToken: (token: string) => void
  changeLimit: (limit: number) => void
  changeCosine: (cosine: number) => void
}

const UserSettingContext = createContext<Setting>({
  openAIKey: null,
  githubToken: null,
  limit: 5,
  cosine: 0.8,
  changeKey: (key: string) => {},
  changeToken: (token: string) => {},
  changeLimit: (limit: number) => {},
  changeCosine: (cosine: number) => {},
})

const SettingProvider = ({ children, info }: { children: React.ReactNode; info: UserInfo }) => {
  const [openAIKey, setOpenAIKey] = useState<string | null>(info?.openAIKey ?? null)
  const [githubToken, setGithubToken] = useState<string | null>(info?.githubToken ?? null)
  const [limit, setLimit] = useState<number>(info?.limit ?? 5)
  const [cosine, setCosine] = useState<number>(info?.cosine ?? 0.8)

  const changeKey = (key: string) => {
    setOpenAIKey(key)
  }

  const changeToken = (token: string) => {
    setGithubToken(token)
  }

  const changeLimit = (limit: number) => {
    setLimit(limit)
  }

  const changeCosine = (cosine: number) => {
    setCosine(cosine)
  }

  return (
    <UserSettingContext.Provider
      value={{
        openAIKey: openAIKey,
        githubToken: githubToken,
        limit: limit,
        cosine: cosine,
        changeKey,
        changeToken,
        changeLimit,
        changeCosine,
      }}
    >
      {children}
    </UserSettingContext.Provider>
  )
}

const useSetting = () => {
  const context = useContext(UserSettingContext)

  return context
}

export { SettingProvider, useSetting }
