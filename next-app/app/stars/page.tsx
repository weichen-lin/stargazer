import clsx from 'clsx'

export const dynamic = 'force-dynamic'

export default async function Dashboard() {
  return (
    <div
      className={clsx(
        'bg-blue-50/70 dark:bg-slate-900 flex flex-col items-center',
        'w-full lg:h-full md:flex-row md:flex-wrap md:items-center md:justify-center',
        'lg:py-0 lg:p-2 lg:gap-1 xl:gap-3 2xl:gap-x-12 2xl:gap-y-1',
        'from-blue-100/40 to-blue-200/60 bg-gradient-to-br dark:from-slate-500 dark:to-slate-100',
      )}
    >
      <div
        className={clsx(
          'flex flex-col justify-between items-center backdrop-blur-md backdrop-opacity-60',
          'w-[380px] h-[380px] xl:w-[330px] 2xl:w-[380px] xl:h-[330px] 2xl:h-[380px]',
        )}
      ></div>
    </div>
  )
}
