'use client'

import { motion } from 'framer-motion'
import { ISuggestion } from '@/actions'
import {
  ExclamationTriangleIcon,
  GearIcon,
  GitHubLogoIcon,
  CrossCircledIcon,
  OpenInNewWindowIcon,
} from '@radix-ui/react-icons'
import { useRepoDetail } from '@/hooks/util'

export const Empty = () => {
  return (
    <div className='py-3 flex flex-col lg:flex-row gap-y-1 flex-wrap w-full gap-4 bg-slate-300/30 p-3 lg:w-2/3'>
      <div className='flex items-center gap-x-4'>
        <ExclamationTriangleIcon className='w-5 h-5 text-yellow-500' />
        <p className='text-slate-700'>No data found</p>
      </div>
      <div className='text-slate-500 text-left'>
        Please click on the symbol in the upper right corner
        <GearIcon className='w-5 h-5 inline-block mx-1 mb-1' />
        to adjust your cosine similarity threshold.
      </div>
    </div>
  )
}

export const Error = () => {
  return (
    <div className='py-3 flex flex-col lg:flex-row gap-y-1 flex-wrap w-full gap-4 bg-slate-300/30 p-3 lg:w-2/3'>
      <div className='flex items-center gap-x-4'>
        <CrossCircledIcon className='w-5 h-5 text-red-500' />
        <p className='text-slate-700'>Service not available now</p>
      </div>
    </div>
  )
}

export const HaveSuggestions = (props: { suggestions: ISuggestion[] }) => {
  const { suggestions } = props

  return suggestions.length > 0 ? (
    <div className='w-full lg:w-2/3 lg:mx-auto flex flex-col gap-y-2 rounded-md bg-slate-300/30'>
      <p className='text-slate-700 p-3 px-5'>Stargazer found {suggestions.length} results for you</p>
      <div className='py-3 flex flex-col md:grid md:grid-cols-2 gap-y-5 flex-wrap w-full gap-4 p-3'>
        {suggestions.map((e, i) => (
          <Suggestion key={i} {...e} index={i} />
        ))}
      </div>
    </div>
  ) : (
    <Empty />
  )
}

const Suggestion = (props: ISuggestion & { index: number }) => {
  const { repo_id, avatar_url, full_name, html_url, index } = props
  const { setOpen, setRepoID } = useRepoDetail()

  return (
    <motion.div
      initial={{ opacity: 0, y: -100 }}
      animate={{ opacity: 1, y: 0, transition: { delay: index * 0.05 } }}
      className='rounded-md drop-shadow-md shadow-md w-full flex bg-white border-[1px] border-slate-100/40 p-2 flex-col gap-y-2 max-w-[360px]'
    >
      <div className='flex justify-between w-full items-center'>
        <div className='flex gap-x-2 items-center'>
          <div className='w-10 h-10'>
            <img src={avatar_url} alt='stargazer' className='rounded-full w-full h-full' />
          </div>
          <div className='max-w-[220px] truncate'>{full_name}</div>
        </div>
        <div className='flex gap-x-2 items-center'>
          <OpenInNewWindowIcon
            className='cursor-pointer'
            onClick={() => {
              setOpen(true)
              setRepoID(repo_id)
            }}
          />
          <a className='p-2 rounded-full hover:bg-slate-300/30' href={html_url} target='_blank'>
            <GitHubLogoIcon />
          </a>
        </div>
      </div>
    </motion.div>
  )
}
