import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { useSearch } from '@/hooks/stars'
import { motion } from 'framer-motion'
import { Input } from '@/components/ui/input'
import { useState, useRef } from 'react'
import { cn } from '@/lib/utils'
import clsx from 'clsx'

export default function Search() {
  const { query, setQuery, open, setOpen, repos } = useSearch()
  const ref = useRef<HTMLInputElement>(null)

  return (
    <Dialog open={open} onOpenChange={e => setOpen(e)}>
      <DialogTrigger>Open</DialogTrigger>
      <DialogContent className='border-0 bg-slate-100 top-[20%]'>
        <div className={cn('flex flex-col gap-y-4')}>
          <Input
            ref={ref}
            value={query}
            onChange={e => {
              setQuery(e.target.value)
            }}
            className={clsx(
              'text-slate-700 rounded-md p-4 bg-white',
              'border-[1px] border-slate-300 focus:shadow-md focus:outline-none focus:ring-0',
              'focus-visible:ring-0 focus:border-slate-500',
            )}
            placeholder='Search for a repository...'
          />
          {repos.length > 0 && (
            <motion.div
              initial={{ opacity: 0, scaleY: 0 }}
              animate={{ opacity: 1, scaleY: 1 }}
              transition={{ duration: 0.3 }}
              className='flex-1 rounded-b-md h-0 origin-top'
            >
              {repos.map((repo, i) => (
                <motion.a
                  key={i}
                  href={repo.html_url}
                  target='_blank'
                  rel='noreferrer'
                  className='flex flex-col items-start gap-y-4 p-2 bg-white rounded-lg hover:shadow-md transition-colors'
                >
                  <div className='flex gap-x-2 items-center'>
                    <img src={repo.avatar} alt={repo.full_name} className='w-4 h-4 rounded-md' />
                    <div className='font-semibold text-blue-700'>{repo.full_name}</div>
                  </div>
                  <div className='flex flex-col pl-6 gap-y-2'>
                    <div className='text-slate-400 text-sm'>{repo.description}</div>
                    <div className='text-slate-700 text-sm'>{repo.summary}</div>
                  </div>
                </motion.a>
              ))}
            </motion.div>
          )}
          {repos.length === 0 && (
            <div className='rounded-b-md origin-top w-full bg-slate-100 text-slate-500/70 p-8 text-center'>
              No recent searches
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
