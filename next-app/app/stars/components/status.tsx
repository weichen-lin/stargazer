import Image from 'next/image'

export function Empty() {
  return (
    <div className='flex flex-col justify-center items-center gap-y-1 w-full h-full'>
      <Image src='/empty.png' alt='Empty' width={380} height={380} />
      <p className='text-md text-gray-600 dark:text-gray-400 px-8 text-center'>
        Discover more with our powerful search feature. Try it out now!
      </p>
    </div>
  )
}

export function Loading() {
  return (
    <div className={`h-full md:grid md:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4 gap-4 pt-4 flex flex-col mx-auto`}>
      {Array.from({ length: 10 }).map((_, index) => (
        <div className='w-[380px] mx-auto h-[170px] animate-pulse bg-slate-200' key={`loading-skeleton-${index}`}></div>
      ))}
    </div>
  )
}
