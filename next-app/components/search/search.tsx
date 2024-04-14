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
import { useState } from 'react'

export default function Search() {
  const { query, setQuery, open, setOpen } = useSearch()
  const [repos, setRepos] = useState([])

  return (
    <Dialog open={open} onOpenChange={e => setOpen(e)}>
      <DialogTrigger>Open</DialogTrigger>
      <DialogContent className='h-[400px] p-0 bg-transparent border-0'>
        <div className='flex flex-col p-[5px] gap-y-1'>
          <Input
            value={query}
            onChange={e => {
              setQuery(e.target.value)
            }}
            className='bg-slate-100 text-slate-900 rounded-md'
            placeholder='Search for a repository...'
          />
          {repos.length > 0 && (
            <motion.div
              initial={{ opacity: 0, scaleY: 0 }}
              animate={{ opacity: 1, scaleY: 1 }}
              transition={{ duration: 0.5 }}
              className='flex-1 bg-slate-100 rounded-md h-0 origin-top'
            ></motion.div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
