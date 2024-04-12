'use server'

import { MobileSidebar, DesktopBar } from '@/components/sidebar/chat'
import Chatter from '@/components/util/chatter'
import { getUserInfo } from '@/actions/neo4j'
import { getServerSession } from 'next-auth'
import { redirect } from 'next/navigation'
import { options } from '@/app/api/auth/[...nextauth]/option'
import { SettingProvider } from '@/components/provider/setting'

export default async function ChatPage() {
  const session = await getServerSession(options)

  if (!session) {
    redirect('/')
  }

  const email = session?.user?.email ?? null

  if (!email) {
    redirect('/')
  }

  const info = await getUserInfo({ email: email })
  if (!info) {
    redirect('/')
  }

  return (
    <div className='flex flex-col w-full h-screen'>
      <SettingProvider info={info}>
        <MobileSidebar />
        <DesktopBar />
        <div className='flex-1 w-full overflow-auto'>
          <Chatter />
        </div>
      </SettingProvider>
    </div>
  )
}
