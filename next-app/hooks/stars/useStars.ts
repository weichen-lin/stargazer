'use client'

import { create } from 'zustand'
import { Option } from '@/components/ui/multiple-selector'
import type { Repository } from '@/actions/neo4j'

interface SelectState {
  selected: Option[]
  setSelected: (selected: Option[]) => void
  isSearching: boolean
  setIsSearching: (isSearching: boolean) => void
  count: number
  setCount: (count: number) => void
  page: number
  setPage: (page: number) => void
  results: Repository[]
  setResults: (results: Repository[]) => void
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

  const search = async (page: number) => {
    if (selected.length === 0) return
    try {
      setIsSearching(true)
      setPage(page)
      const results = await fetch('/api/repos', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          page: page,
          languages: selected.map(({ value }) => value),
          limit: 20,
        }),
      })

      const { total, repos } = await results.json()
      setResults(repos)
      setCount(total)
    } catch (error) {
      console.error(error)
    } finally {
      setIsSearching(false)
    }
  }

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
