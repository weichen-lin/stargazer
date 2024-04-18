'use client'

import { Button } from '@/components/ui/button'
import { SymbolIcon } from '@radix-ui/react-icons'
import { useState } from 'react'
import { syncUserStars } from '@/actions/producer'
import { useSession } from 'next-auth/react'
import { useToast } from '@/components/ui/use-toast'

export default function SyncStars() {
  const [syncing, setSyncing] = useState(false)
  const { data: session } = useSession()
  const { toast } = useToast()

  const getStars = async () => {
    if (session?.user) {
      setSyncing(true)
      const { status, title, message } = await syncUserStars(session?.user?.name ?? '')
      setSyncing(false)
      if (status === 200) {
        toast({
          title,
          description: message,
        })
      } else {
        toast({
          title,
          description: message,
          variant: 'destructive',
        })
      }
    }
  }

  return (
    <Button variant='outline' size='icon' className='ml-4' loading={syncing} onClick={getStars}>
      <SymbolIcon />
    </Button>
  )
}
