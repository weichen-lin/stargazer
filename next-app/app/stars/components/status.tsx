import Image from 'next/image'

export function Empty() {
  return (
    <div className='flex flex-col justify-center items-center gap-y-4 w-full h-full'>
      <Image src='/empty.png' alt='Empty' width={380} height={380} />
      <p className='text-md text-gray-600 dark:text-gray-400'>
        Discover more with our powerful search feature. Try it out now!
      </p>
    </div>
  )
}

export function Loading() {
  return (
    <div className='w-full h-full grid grid-cols-4 gap-4 py-4 items-start'>
      {Array.from({ length: 10 }).map((_, index) => (
        <div className='h-[170px] animate-pulse bg-slate-200' key={`loading-skeleton-${index}`}></div>
      ))}
    </div>
  )
}
