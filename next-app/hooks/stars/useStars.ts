'use client'

import { create } from 'zustand'
import { Option } from '@/components/ui/multiple-selector'
import { useFetch } from '@/hooks/util'
import { IRepoSearchWithLanguage, IRepository } from '@/client/repository'
import { useCallback } from 'react'

interface SelectState {
  selected: Option[]
  setSelected: (selected: Option[]) => void
  isSearching: boolean
  setIsSearching: (isSearching: boolean) => void
  count: number
  setCount: (count: number) => void
  page: number
  setPage: (page: number) => void
  results: IRepository[]
  setResults: (results: IRepository[]) => void
}

const selectState = create<SelectState>(set => ({
  selected: [],
  setSelected: selected => set({ selected }),
  isSearching: false,
  setIsSearching: isSearching => set({ isSearching }),
  count: 0,
  setCount: count => set({ count }),
  page: 1,
  setPage: page => set({ page }),
  results: [],
  setResults: results => set({ results }),
}))

const useStars = () => {
  const { selected, setSelected, count, setCount, page, setPage, isSearching, setIsSearching, results, setResults } =
    selectState()

  const { run } = useFetch<IRepoSearchWithLanguage>({
    initialRun: false,
    config: {
      url: '/repository',
      method: 'GET',
    },
    onSuccess: data => {
      setResults(data.data)
      setCount(data.total)
    },
  })

  const search = useCallback(
    async (page: number) => {
      if (selected.length === 0) return
      try {
        setIsSearching(true)
        setPage(page)
        await run({
          params: { languages: selected.map(e => e.value).join(','), page: page.toString(), limit: '20' },
        })
      } catch (error) {
        console.error(error)
      } finally {
        setIsSearching(false)
      }
    },
    [selected, run, setCount, setIsSearching, setPage],
  )

  return {
    selected,
    setSelected,
    isSearching,
    count,
    results,
    search,
    page,
  }
}

export default useStars
