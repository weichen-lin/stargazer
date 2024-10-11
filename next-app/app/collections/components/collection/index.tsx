import { cn } from '@/lib/utils'
import Collection from './collection'
import { useCollection } from '@/app/collections/hooks'
import Image from 'next/image'

export default function Collections() {
  const { data, total, loading } = useCollection()

  if (loading) {
    return <Loading />
  }

  if (!data || data.length === 0) {
    return <Empty />
  }

  return (
    <div
      className={cn(
        'w-full flex flex-col py-8 h-full md:pl-3',
        'content-start justify-start items-start',
        'flex-1 md:grid md:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4 gap-4 overflow-y-auto',
      )}
    >
      {data.map(collection => (
        <Collection key={collection.id} {...collection} />
      ))}
    </div>
  )
}

const Loading = () => {
  return (
    <div
      className={cn(
        'w-full flex flex-col py-8 h-full md:pl-3',
        'content-start justify-start items-start',
        'flex-1 md:grid md:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4 gap-4',
      )}
    >
      <div className='h-[100px] bg-slate-200 animate-pulse w-full rounded-lg'></div>
      <div className='h-[100px] bg-slate-200 animate-pulse w-full rounded-lg'></div>
      <div className='h-[100px] bg-slate-200 animate-pulse w-full rounded-lg'></div>
      <div className='h-[100px] bg-slate-200 animate-pulse w-full rounded-lg'></div>
    </div>
  )
}

const Empty = () => {
  return (
    <div className='w-full h-full flex flex-col items-center justify-center gap-y-4'>
      <Image src='/empty-collection.png' alt='Empty' width={570} height={570} />
      <span className='text-slate-400'>No collections found</span>
    </div>
  )
}
