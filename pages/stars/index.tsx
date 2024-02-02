'use client'

import { ArrangeSetting } from '@/components/tab'
import { useArrangement } from '@/hooks/stars'
import { GridRepo } from '@/components/repo'
import { motion } from 'framer-motion'

export default function Stars() {
  const { arrangement, toggleArrangement } = useArrangement()

  return (
    <div className='flex flex-col gap-y-12 w-full'>
      <div className='flex items-center justify-between w-full'>
        <motion.h1 initial={{ x: 80 }} animate={{ x: 0 }} className='text-4xl font-semibold'>
          My Stars
        </motion.h1>
        <ArrangeSetting arrangement={arrangement} toggle={toggleArrangement} />
      </div>
      <div className='w-full flex flex-wrap gap-8 flex-1 pb-4 items-center justify-start'>
        {Array.from({ length: 10 }).map((_, i) => (
          <GridRepo key={i} {...githubStar} index={i} />
        ))}
      </div>
    </div>
  )
}

const githubStar = {
  id: 107505869,
  full_name: 'firecracker-microvm/firecracker',
  owner: {
    avatar_url: 'https://avatars.githubusercontent.com/u/44477506?v=4',
  },
  html_url: 'https://github.com/firecracker-microvm/firecracker',
  description: 'Secure and fast microVMs for serverless computing.',
  updated_at: '2024-02-02T04:50:52Z',
  homepage: 'http://firecracker-microvm.io',
  stargazers_count: 23457,
  language: 'Rust',
}
