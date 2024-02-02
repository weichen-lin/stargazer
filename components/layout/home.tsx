import Sidebar from '@/components/sidebar'
import Image from 'next/image'

export default async function HomeLayout(props: { username: string; children: React.ReactNode; path: string }) {
  const { username, path, children } = props
  return (
    <div className='bg-[#f2f0e8] h-screen w-screen flex'>
      <div className='pl-4 flex flex-col gap-y-8 w-[240px]'>
        <div className='py-6 w-full flex items-center justify-center border-b-[1px] border-slate-500/10'>
          <Image
            src='/icon.png'
            alt='Logo'
            width={50}
            height={50}
            className='rounded-full cursor-pointer hover:opacity-75'
          />
        </div>
        <Sidebar path={path} username={username} />
        <div className=''></div>
      </div>
      <div className='flex-1 h-full overflow-y-auto'>{children}</div>
    </div>
  )
}
