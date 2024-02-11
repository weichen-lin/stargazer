'use server'

import { MobileSidebar, DesktopBar } from '@/components/sidebar/chat'
import Chatter from '@/components/util/chatter'

export default async function ChatPage() {
  return (
    <div className='flex flex-col w-full h-screen'>
      <MobileSidebar />
      <DesktopBar />
      <div className='flex-1 w-full overflow-auto'>
        <Chatter />
      </div>
    </div>
  )
}
