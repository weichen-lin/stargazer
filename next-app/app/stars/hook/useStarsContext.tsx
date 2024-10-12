'use client'

import { createContext, useContext, useState, useCallback } from 'react'
import { IRepoSearchWithLanguage, IRepositoryWithCollection } from '@/client/repository'
import { useFetch, IRunProps } from '@/hooks/util'
import { Option } from '@/components/ui/multiple-selector'

interface IStarsContext {
  selectLanguages: Option[]
  setSelectLanguages: (languages: Option[]) => void
  selectRepos: number[]
  setSelectRepos: (repos: number[]) => void
  page: number
  setPage: (page: number) => void
  total: number
  setTotal: (total: number) => void
  data: IRepositoryWithCollection[]
  setData: (data: IRepositoryWithCollection[]) => void
  open: boolean
  setOpen: (open: boolean) => void
  isSearching: boolean
  search: (page: number) => Promise<void>
}

export const StarsContext = createContext<IStarsContext | null>(null)

export const useStarsContext = (): IStarsContext => {
  const context = useContext(StarsContext)
  if (context === null) {
    throw new Error('StarsContext must be used within a CollectionProvider')
  }

  return context
}

export const StarsProvider = ({ children }: { children: React.ReactNode }) => {
  const [selectLanguages, setSelectLanguages] = useState<Option[]>([])
  const [selectRepos, setSelectRepos] = useState<number[]>([])
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [data, setData] = useState<IRepositoryWithCollection[]>([])
  const [open, setOpen] = useState(false)
  const [isSearching, setIsSearching] = useState(false)

  const { run } = useFetch<IRepoSearchWithLanguage>({
    initialRun: false,
    config: {
      url: '/repository',
      method: 'GET',
    },
    onSuccess: data => {
      setData(data.data)
      setTotal(data.total)
    },
  })

  const search = useCallback(
    async (page: number) => {
      if (selectLanguages.length === 0) return
      try {
        setIsSearching(true)
        setPage(page)
        await run({
          params: { languages: selectLanguages.map(e => e.value).join(','), page: page.toString(), limit: '20' },
        })
      } catch (error) {
        console.error(error)
      } finally {
        setIsSearching(false)
      }
    },
    [selectLanguages, run, setTotal, setIsSearching, setPage],
  )

  return (
    <StarsContext.Provider
      value={{
        selectLanguages,
        setSelectLanguages,
        selectRepos,
        setSelectRepos,
        page,
        setPage,
        total,
        setTotal,
        data,
        setData,
        open,
        setOpen,
        isSearching,
        search,
      }}
    >
      {children}
    </StarsContext.Provider>
  )
}
