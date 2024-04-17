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
import { ISuggestion } from '@/actions'

export default function Search() {
  const { query, setQuery, open, setOpen, repos, ref } = useSearch()

  return (
    <Dialog open={open} onOpenChange={e => setOpen(e)}>
      <DialogTrigger>Open</DialogTrigger>
      <DialogContent className='border-0 bg-slate-100 top-[30%] max-h-[500px]'>
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
            <div className='origin-top overflow-y-scroll h-[300px]'>
              {repos.map((repo, i) => (
                <Repo {...repo} query={query} />
              ))}
            </div>
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

function FullTextIndex(s: string, target: string): number[] | null {
  let left = 0
  let right = 0
  const result: number[] = []

  while (right < s.length && left < target.length) {
    if (s[right].toLowerCase() === target[left]) {
      result.push(right)
      left++
    }
    right++
  }

  return result.length === target.length ? result : null
}

const Repo = (props: ISuggestion & { query: string }) => {
  const { avatar_url, full_name, description, readme_summary, query } = props

  const full_name_index = FullTextIndex(full_name, query)
  const description_index = FullTextIndex(description ?? '', query)
  const summary_index = FullTextIndex(readme_summary, query)

  return (
    <div className='flex flex-col items-start gap-y-2 p-2 bg-white rounded-lg hover:shadow-md transition-colors'>
      <div className='flex gap-x-2 items-center'>
        <img src={avatar_url} alt={full_name} className='w-4 h-4 rounded-md' />
        <div className='font-semibold text-blue-700 line-clamp-1'>
          {full_name_index ? (
            <>
              {full_name.split('').map((char, i) => {
                if (full_name_index.includes(i)) {
                  return <span className='underline decoration-blue-300'>{char}</span>
                }
                return <>{char}</>
              })}
            </>
          ) : (
            full_name
          )}
        </div>
      </div>
      <div className='flex flex-col pl-6 gap-y-2'>
        <div className='text-slate-400 text-sm line-clamp-2'>
          {description_index && description ? (
            <>
              {description.split('').map((char, i) => {
                if (description_index.includes(i)) {
                  return <span className='underline decoration-blue-500'>{char}</span>
                }
                return <>{char}</>
              })}
            </>
          ) : (
            description ?? ''
          )}
        </div>
        <div className='text-slate-700 text-sm line-clamp-3'>
          {summary_index ? (
            <>
              {readme_summary.split('').map((char, i) => {
                if (summary_index.includes(i)) {
                  return <span className='underline decoration-blue-300'>{char}</span>
                }
                return <>{char}</>
              })}
            </>
          ) : (
            readme_summary
          )}
        </div>
      </div>
    </div>
  )
}
