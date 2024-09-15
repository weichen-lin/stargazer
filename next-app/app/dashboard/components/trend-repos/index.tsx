'use client'

import { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { StarFilledIcon, EyeOpenIcon, RocketIcon } from '@radix-ui/react-icons'
import { PieChart as PieChartIcon, Plus } from 'lucide-react'
import { SortKey, IRepository } from '@/client/repository'
import { useFetch } from '@/hooks/util'
import TrendRepository, { DateRange, sinceMap } from './trending-repo'
import LanguageSelector from './language-selector'
import { Badge } from '@/components/ui/badge'

const fakeData = [
  {
    repo_name: 'ecapture',
    owner_name: 'gojue',
    html_url: 'https://github.com/gojue/ecapture',
    description:
      'Capturing SSL/TLS plaintext without a CA certificate using eBPF. Supported on Linux/Android kernels for amd64/arm64.',
    stargazers_count: 10674,
    language: 'C',
    get_stars: 514,
  },
]

function capitalizeFirstLetter(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1)
}

export default function TrendRepos() {
  const [since, setSince] = useState<DateRange>('daily')
  // const { isLoading, run } = useFetch<IRepository[]>({
  //   initialRun: false,
  //   config: {
  //     url: '/repository/sort',
  //     method: 'GET',
  //     params: {
  //       key: sortKey,
  //     },
  //   },
  // })

  return (
    <Card className='flex flex-col h-[320px] w-full max-w-[380px] md:max-w-none'>
      <CardHeader className='items-start pb-0 gap-y-0'>
        <CardTitle className='text-xl flex gap-x-4 justify-between w-full'>
          <div className='flex items-center gap-x-3'>
            <div>Trending Repositories</div>
            <div className='mb-1'>
              <Badge variant='secondary'>from GitHub</Badge>
            </div>
          </div>
        </CardTitle>
        <CardDescription className='flex gap-x-3'>
          <LanguageSelector />
          <Select
            value={since}
            onValueChange={e => {
              setSince(e as DateRange)
            }}
          >
            <SelectTrigger className='w-[180px]'>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {Object.entries(sinceMap).map(([key, value]) => (
                <SelectItem value={key}>{capitalizeFirstLetter(value)}</SelectItem>
              ))}
            </SelectContent>
          </Select>
        </CardDescription>
      </CardHeader>
      <CardContent className='flex-1 pb-0 overflow-y-auto flex flex-col gap-y-2 py-4'>
        {fakeData.map(repo => (
          <TrendRepository key={`trend_repo_${repo.owner_name}_${repo.repo_name}`} {...repo} date_range={since} />
        ))}
        {/* {!isLoading && data && data.length > 0 && data.map(repo => <Repo key={`repo_${repo.repo_id}`} {...repo} />)}
        {!isLoading && data && data.length === 0 && <EmptyContent />} */}
        {false && <Loading />}
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
