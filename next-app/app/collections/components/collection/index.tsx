import { cn } from '@/lib/utils'
import Collection from './collection'
import { useCollection } from '@/app/collections/hooks'

export default function Collections() {
  const { data, total, loading } = useCollection()

  return (
    <div
      className={cn(
        'w-full flex flex-col justify-start items-start py-8 h-full md:pl-3 content-start',
        'flex-1 md:grid md:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4 gap-4 overflow-y-auto',
      )}
    >
      {loading && <Loading />}
      {!loading && data?.length === 0 && <Empty />}
      {!loading && data.map(collection => <Collection key={collection.id} {...collection} />)}
    </div>
  )
}

const Loading = () => {
  return (
    <>
      <div className='h-[100px] bg-slate-200 animate-pulse w-full rounded-lg'></div>
      <div className='h-[100px] bg-slate-200 animate-pulse w-full rounded-lg'></div>
      <div className='h-[100px] bg-slate-200 animate-pulse w-full rounded-lg'></div>
      <div className='h-[100px] bg-slate-200 animate-pulse w-full rounded-lg'></div>
    </>
  )
}

const Empty = () => {
  return (
    <div className='w-full flex items-center justify-center'>
      <span className='text-slate-400'>No collections found</span>
    </div>
  )
}
