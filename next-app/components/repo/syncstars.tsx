'use client'

import { Button } from '@/components/ui/button'
import { SymbolIcon } from '@radix-ui/react-icons'
import { useState } from 'react'
import { syncUserStars } from '@/actions/producer'
import { useToast } from '@/components/ui/use-toast'
import { useUser } from '@/context/user'

export default function SyncStars() {
  const [syncing, setSyncing] = useState(false)

  const { name } = useUser()
  const { toast } = useToast()

  const getStars = async () => {
    setSyncing(true)
    const { status, title, message } = await syncUserStars(name)
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

  return (
    <Button variant='outline' size='icon' className='ml-4' loading={syncing} onClick={getStars}>
      <SymbolIcon />
    </Button>
  )
}
