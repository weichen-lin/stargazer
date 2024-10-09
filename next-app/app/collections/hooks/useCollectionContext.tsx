'use client'

import { createContext, useContext } from 'react'
import { ICollection } from '@/client/collection'
import { useFetch } from '@/hooks/util'

interface ICollectionContext extends ICollection {
  setCollection: (collection: ICollection) => void
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
  collection,
}: {
  children: React.ReactNode
  collection: ICollection
}) => {
  const setCollection = (collection: ICollection) => {}

  return <CollectionContext.Provider value={{ ...collection, setCollection }}>{children}</CollectionContext.Provider>
}
