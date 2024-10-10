'use client'

import { createContext, useContext, useState } from 'react'
import { ICollection } from '@/client/collection'
import { useFetch, IRunProps } from '@/hooks/util'

interface ICollectionContext {
  collection: ICollection
  update: (collection: ICollection) => void
  isUpdate: boolean
  isSearch: boolean
  repos: any[]
  page: number
  total: number
  getCollectRepos: (params: IRunProps) => void
}

export const CollectionContext = createContext<ICollectionContext | null>(null)

export const useCollectionContext = (): ICollectionContext => {
  const context = useContext(CollectionContext)
  if (context === null) {
    throw new Error('useCollectionContext must be used within a CollectionProvider')
  }

  return context
}

export const CollectionProvider = ({
  children,
  initCollection,
}: {
  children: React.ReactNode
  initCollection: ICollection
}) => {
  const [collection, setCollection] = useState(initCollection)
  const [selectRepos, setSelectRepos] = useState([])
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [repos, setRepos] = useState([])

  const { run, isLoading } = useFetch<ICollection>({
    config: {
      url: `/collection/${collection.id}`,
      method: 'PATCH',
    },
    initialRun: false,
    onSuccess: data => setCollection(data),
  })

  const { run: getCollectRepos, isLoading: getCollectReposLoading } = useFetch<ICollection[]>({
    config: {
      url: `/collection/repos/${collection.id}`,
      method: 'GET',
      params: {
        page,
        limit: 20,
      },
    },
    initialRun: true,
  })

  const update = async (collection: ICollection) => {
    run({ payload: { name: collection.name, description: collection.description, is_public: collection.is_public } })
  }

  return (
    <CollectionContext.Provider
      value={{
        collection,
        isUpdate: isLoading,
        isSearch: getCollectReposLoading,
        update,
        repos,
        total,
        page,
        getCollectRepos,
      }}
    >
      {children}
    </CollectionContext.Provider>
  )
}
