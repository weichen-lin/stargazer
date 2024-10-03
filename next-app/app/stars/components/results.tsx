'use client'

import { useStars } from '@/hooks/stars'
import { Loading, Empty } from './status'
import Repo from './repo'

export default function Results() {
  const { isLoading, results } = useStars()

  if (isLoading) {
    return <Loading />
  }

  if (!results || results.length === 0) {
    return <Empty />
  }

  return (
    <div className={`w-full h-full md:grid md:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4 gap-4 py-4 flex flex-col`}>
      {results.map((item, index) => (
        <Repo key={item.repo_id} {...item} index={index} />
      ))}
    </div>
  )
}
