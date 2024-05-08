import clsx from 'clsx'

const TotalRepositoriesLoading = () => {
  return (
    <div
      className={clsx(
        'border-[2px] border-slate-300/70 rounded-lg  w-full h-[168px] flex flex-col items-center justify-center gap-y-4',
        'bg-white drop-shadow-lg dark:bg-slate-300 dark:border-slate-800 dark:text-white',
      )}
    >
      <div className='w-12 text-center text-6xl h-8 bg-slate-300 animate-pulse rounded-md'></div>
      <div className='w-36 text-center text-2xl h-12 bg-slate-300 animate-pulse rounded-md'></div>
    </div>
  )
}

const ReposLoading = ({ title }: { title: string }) => {
  return (
    <div
      className={clsx(
        'border-[2px] border-slate-300/70 rounded-lg flex flex-col items-center justify-start gap-y-2',
        'bg-white drop-shadow-lg dark:bg-slate-300 dark:border-slate-800 dark:text-white',
        'w-[380px] h-[380px]',
      )}
    >
      <div className='w-full text-left px-4 py-3 text-xl text-slate-500'>{title}</div>
      {Array.from({ length: 5 }).map((_, index) => (
        <div
          className='w-[90%] mx-auto h-12 animate-pulse bg-slate-300'
          key={`loading-skeleton-${title}-${index}`}
        ></div>
      ))}
    </div>
  )
}

export { TotalRepositoriesLoading, ReposLoading }
