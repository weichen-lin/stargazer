'use client'

import { Button } from '@/components/ui/button'
import { SymbolIcon } from '@radix-ui/react-icons'
import { useState } from 'react'
import { syncUserStars } from '@/actions/producer'
import { useSession } from 'next-auth/react'
import { toast } from 'react-toastify'

export default function SyncStars() {
  const [syncing, setSyncing] = useState(false)
  const { data: session } = useSession()

  const getStars = async () => {
    if (session?.user) {
      setSyncing(true)
      const { status, message } = await syncUserStars(session?.user?.name ?? '')
      setSyncing(false)

      if (status === 200) {
        toast.success(message, { position: 'bottom-center' })
      } else {
        toast.error(message, { position: 'bottom-center' })
      }
    }
  }

  return (
    <Button variant='outline' size='icon' className='ml-4' loading={syncing} onClick={getStars}>
      <SymbolIcon />
    </Button>
  )
}
