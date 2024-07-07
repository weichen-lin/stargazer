import clsx from 'clsx'

export default function Loading() {
  return (
    <div
      className={clsx(
        'flex flex-col items-center dark:bg-black',
        'w-full flex-1 overflow-y-auto px-[2.5%] py-8 gap-y-3',
        'lg:grid lg:grid-cols-2 lg:gap-x-3 lg:gap-y-3 xl:grid-cols-3 xl:gap-x-3 xl:gap-y-3',
        'lg:content-around',
      )}
    >
      {/* <div className='w-full bg-slate-300/40 dark:bg-slate-100/30 h-[100px] animate-pulse rounded-md'></div>
      <div className='w-full bg-slate-300/40 dark:bg-slate-100/30 h-[100px] animate-pulse rounded-md'></div>
      <div className='w-full bg-slate-300/40 dark:bg-slate-100/30 h-[100px] animate-pulse rounded-md'></div>
      <div className='w-full bg-slate-300/40 dark:bg-slate-100/30 h-[100px] animate-pulse rounded-md'></div>
      <div className='w-full bg-slate-300/40 dark:bg-slate-100/30 h-[100px] animate-pulse rounded-md'></div>
      <div className='w-full bg-slate-300/40 dark:bg-slate-100/30 h-[100px] animate-pulse rounded-md'></div> */}
      {/* <div className='w-full bg-slate-300/40 dark:bg-slate-100/30 animate-pulse rounded-md h-[300px]'></div>
      <div className='w-full bg-slate-300/40 dark:bg-slate-100/30 animate-pulse rounded-md h-[300px]'></div>
      <div className='w-full bg-slate-300/40 dark:bg-slate-100/30 animate-pulse rounded-md h-[300px]'></div>
      <div className='w-full bg-slate-300/40 dark:bg-slate-100/30 animate-pulse rounded-md h-[300px]'></div>
      <div className='w-full bg-slate-300/40 dark:bg-slate-100/30 animate-pulse rounded-md h-[300px]'></div>
      <div className='w-full bg-slate-300/40 dark:bg-slate-100/30 animate-pulse rounded-md h-[300px]'></div> */}
    </div>
  )
}
