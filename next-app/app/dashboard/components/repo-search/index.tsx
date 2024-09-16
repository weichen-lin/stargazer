'use client'

import { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { StarFilledIcon, EyeOpenIcon, RocketIcon } from '@radix-ui/react-icons'
import { PieChart as PieChartIcon, Plus } from 'lucide-react'
import { SortKey, IRepository } from '@/client/repository'
import { useFetch } from '@/hooks/util'
import Repo from './repo'
import { Badge } from '@/components/ui/badge'

const CardInfoMap: {
  [key in SortKey]: { title: string; short: string; description: string; icon: JSX.Element }
} = {
  created_at: {
    title: 'Newest Additions',
    short: 'Recent',
    description: "Discover the latest repositories you've starred.",
    icon: <RocketIcon />,
  },
  stargazers_count: {
    title: 'Most Popular',
    short: 'Starred',
    description: 'Explore the repositories with the most stars.',
    icon: <StarFilledIcon />,
  },
  watchers_count: {
    title: 'Trending Now',
    short: 'Watched',
    description: "Check out the repositories that everyone's talking about.",
    icon: <EyeOpenIcon />,
  },
}

export default function RepoSearch() {
  const [sortKey, setSortKey] = useState<SortKey>('created_at')
  const { data, isLoading, run } = useFetch<IRepository[]>({
    initialRun: false,
    config: {
      url: '/repository/sort',
      method: 'GET',
      params: {
        key: sortKey,
      },
    },
  })

  useEffect(() => {
    run({ params: { key: sortKey } })
  }, [sortKey])

  return (
    <Card className='flex flex-col h-[320px] w-full max-w-[380px] md:max-w-none'>
      <CardHeader className='items-start pb-0 gap-y-0'>
        <CardTitle className='text-xl flex gap-x-4 justify-between w-full'>
          <div className='flex items-center gap-x-3'>
            <div>{CardInfoMap[sortKey].title}</div>
            <div>
              <Badge variant='secondary'>from your stars</Badge>
            </div>
          </div>
          <Select
            value={sortKey}
            onValueChange={e => {
              setSortKey(e as SortKey)
            }}
            disabled={isLoading}
          >
            <SelectTrigger className='w-[180px]'>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                {Object.entries(CardInfoMap).map(([key, { icon, short }]) => (
                  <SelectItem key={key} value={key}>
                    <div className='flex gap-x-3 items-center'>
                      {icon}
                      <span>{short}</span>
                    </div>
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        </CardTitle>
        <CardDescription> {CardInfoMap[sortKey].description}</CardDescription>
      </CardHeader>
      <CardContent className='flex-1 pb-0 overflow-y-scroll flex flex-col gap-y-2 py-4'>
        {!isLoading && data && data.length > 0 && data.map(repo => <Repo key={`repo_${repo.repo_id}`} {...repo} />)}
        {!isLoading && data && data.length === 0 && <EmptyContent />}
        {isLoading && <Loading />}
      </CardContent>
    </Card>
  )
}

const Loading = () => {
  return (
    <div className='flex flex-col gap-y-2'>
      <div className='w-full rounded-lg h-12 bg-slate-200 animate-pulse'></div>
      <div className='w-full rounded-lg h-12 bg-slate-200 animate-pulse'></div>
      <div className='w-full rounded-lg h-12 bg-slate-200 animate-pulse'></div>
    </div>
  )
}

const EmptyContent = () => {
  return (
    <div className='w-full flex flex-col items-center justify-center mt-4'>
      <div className='w-32 h-32 relative'>
        <PieChartIcon className='w-full h-full text-gray-200' />
        <div className='absolute inset-0 flex items-center justify-center'>
          <Plus className='w-8 h-8 text-gray-400' />
        </div>
      </div>
      <p className='text-center text-gray-500 mb-4'>No data yet.</p>
    </div>
  )
}
