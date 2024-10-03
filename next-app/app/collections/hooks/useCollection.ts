import { create } from 'zustand'
import { useFetch } from '@/hooks/util'
import { ICollection } from '@/client/collection'

interface ICollectionAtom {
  isFetching: boolean
  setIsFetching: (e: boolean) => void
  total: number
  data: ICollection[]
  setTotal: (total: number) => void
  setData: (data: ICollection[]) => void
}

const collectionAtom = create<ICollectionAtom>(set => ({
  isFetching: false,
  setIsFetching: isFetching => set({ isFetching }),
  total: 0,
  data: [],
  setTotal: total => set({ total }),
  setData: data => set({ data }),
}))

export default function useCollection() {
  const { isFetching, setIsFetching, total, data, setTotal, setData } = collectionAtom()
  const { isLoading, run } = useFetch<{
    total: number
    data: ICollection[]
  }>({
    initialRun: true,
    config: {
      url: '/collection',
      method: 'GET',
      params: {
        limit: 20,
        page: 1,
      },
    },
    onSuccess: data => {
      setData(data.data)
      setTotal(data.total)
    },
  })

  const loading = isFetching || isLoading

  return {
    loading,
    setIsFetching,
    total,
    data,
    setTotal,
    setData,
  }
}
