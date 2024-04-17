import { Dialog, DialogContent, DialogTrigger } from '@/components/ui/dialog'
import { useSearch } from '@/hooks/stars'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import clsx from 'clsx'
import FullTextSearchResult from './searchResult'

export default function Search() {
  const { query, setQuery, open, setOpen, repos, ref, loading } = useSearch()

  return (
    <Dialog open={open} onOpenChange={e => setOpen(e)}>
      <DialogTrigger>
        <div
          className={clsx(
            'flex justify-between gap-x-2 items-center',
            'w-[120px] border-[1px] border-slate-300/40 py-1 px-2 rounded-lg',
            'hover:bg-slate-300/40 group',
          )}
        >
          <div className='text-sm text-slate-500 group-hover:text-slate-700'>Search...</div>
          <div className='text-xs tracking-widest opacity-60 bg-slate-300/40 px-2 py-1 rounded-md group-hover:bg-white group-hover:opacity-80'>
            ⌘K
          </div>
        </div>
      </DialogTrigger>
      <DialogContent className='border-0 bg-slate-100 top-[30%] max-h-[500px]'>
        <div className={cn('flex flex-col gap-y-4')}>
          <Input
            ref={ref}
            value={query}
            disabled={loading}
            onChange={e => {
              setQuery(e.target.value)
            }}
            className={clsx(
              'text-slate-700 rounded-md p-4 bg-white w-[95%] mx-auto',
              'border-[1px] border-slate-300 focus:shadow-md focus:outline-none focus:ring-0',
              'focus-visible:ring-0 focus:border-slate-500',
            )}
            placeholder='Search for a repository...'
          />
          {repos.length > 0 && (
            <div className='origin-top overflow-y-auto h-[400px] flex flex-col gap-y-4 relative overflow-visible pb-12'>
              {repos.map((repo, i) => (
                <FullTextSearchResult {...repo} query={query} key={`search_result_${i}`} />
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
