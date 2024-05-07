import clsx from 'clsx'

export default function TotalRepositories() {
  return (
    <div
      className={clsx(
        'border-[2px] border-slate-300/70 rounded-lg  w-full h-[168px] flex flex-col items-center justify-center gap-y-4',
        'bg-white drop-shadow-lg dark:bg-slate-300 dark:border-slate-800 dark:text-white',
      )}
    >
      <div className='w-full text-center text-6xl text-slate-500'>384</div>
      <div className='w-full text-center text-2xl text-slate-500'>Total Repositories</div>
    </div>
  )
}
