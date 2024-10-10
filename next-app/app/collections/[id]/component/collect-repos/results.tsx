'use client'

import { Loading, Empty } from './status'
import Repo from './repo'
import { useCollectionContext } from '@/app/collections/hooks/useCollectionContext'

export default function Results() {
  const { isUpdate, isSearch, repos } = useCollectionContext()

  if (isSearch) {
    return <Loading />
  }

  if (!repos || repos.length === 0) {
    return <Empty />
  }

  return (
    <div className='w-full h-full md:grid md:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4 gap-4 py-4 flex flex-col content-start'>
      {repos.map((item, index) => (
        <Repo key={item.repo_id} {...item} index={index} />
      ))}
    </div>
  )
}
