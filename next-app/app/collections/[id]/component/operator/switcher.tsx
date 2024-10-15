'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Lock, Unlock } from 'lucide-react'
import { cn } from '@/lib/utils'
import { useCollectionContext } from '@/app/collections/hooks/useCollectionContext'
import { Badge } from '@/components/ui/badge'
import { useFetch } from '@/hooks/util'
import { ICollection } from '@/client/collection'

export default function Switcher() {
  const { collection, isSearch, update } = useCollectionContext()
  const [isChecked, setIsChecked] = useState(collection.is_public)
  const { run, isLoading } = useFetch<ICollection>({
    config: {
      url: `/collection/${collection.id}`,
      method: 'PATCH',
    },
    initialRun: false,
    onSuccess: data => {
      update(data)
    },
  })

  const handleToggle = () => {
    setIsChecked(!isChecked)
    run({ payload: { is_public: !isChecked } })
  }

  const loading = isLoading || isSearch

  return (
    <div className='flex gap-x-3 items-center'>
      <button
        onClick={handleToggle}
        className={cn(
          'relative inline-flex h-4 w-10 items-center rounded-full transition-colors mt-[2px]',
          isChecked ? 'bg-green-500' : 'bg-red-500',
          loading && 'cursor-not-allowed opacity-60',
        )}
        disabled={loading}
      >
        <motion.div
          className={cn(
            'inline-flex h-6 w-6 items-center justify-center rounded-full bg-white shadow-lg border-[1px]',
            isChecked ? 'border-green-500' : 'border-red-500',
          )}
          animate={{
            x: isChecked ? 20 : -2,
            rotate: isChecked ? 360 : 0,
          }}
        >
          {isChecked ? <Unlock className='h-4 w-4 text-green-500' /> : <Lock className='h-4 w-4 text-red-500' />}
        </motion.div>
      </button>
      <h2 className='font-light'>{collection.is_public ? <Badge>Public</Badge> : <Badge>Private</Badge>}</h2>
    </div>
  )
}
