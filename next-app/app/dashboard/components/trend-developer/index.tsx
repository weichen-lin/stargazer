'use client'

import { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { PieChart as PieChartIcon, Plus } from 'lucide-react'
import { useFetch } from '@/hooks/util'
import { Badge } from '@/components/ui/badge'
import { DateRange, dateRangeMap, ITrendDeveloper } from '@/client/trends/type'
import TrendDeveloper from './trending-developer'

function capitalizeFirstLetter(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1)
}

export default function TrendDevelopers() {
  const [since, setSince] = useState<DateRange>('daily')
  const { isLoading, run, data } = useFetch<ITrendDeveloper[]>({
    initialRun: false,
    config: {
      url: '/trending/developer',
      method: 'GET',
      params: {
        since,
      },
    },
  })

  useEffect(() => {
    run({ params: { since } })
  }, [since])

  return (
    <Card className='flex flex-col h-[320px] w-full max-w-[380px] md:max-w-none'>
      <CardHeader className='items-start pb-0 gap-y-0'>
        <CardTitle className='text-xl flex gap-x-4 justify-between w-full'>
          <div className='flex items-center gap-x-3'>
            <div>Trending Developers</div>
            <div className='mb-1'>
              <Badge variant='secondary'>from GitHub</Badge>
            </div>
          </div>
        </CardTitle>
        <CardDescription className='flex gap-x-3'>
          <Select
            value={since}
            onValueChange={e => {
              setSince(e as DateRange)
            }}
            disabled={isLoading}
          >
            <SelectTrigger className='w-[180px]'>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {Object.entries(dateRangeMap).map(([key, value]) => (
                <SelectItem value={key} key={`select-item-${key}`}>
                  {capitalizeFirstLetter(value)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </CardDescription>
      </CardHeader>
      <CardContent className='flex-1 pb-0 overflow-y-auto flex flex-col gap-y-2 py-4'>
        {!isLoading &&
          data &&
          data.length > 0 &&
          data.map(developer => (
            <TrendDeveloper key={`trend_developer_${developer.name}_${developer.repo_name}`} {...developer} />
          ))}
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
