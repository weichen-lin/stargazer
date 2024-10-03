'use client'

import { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { PieChart as PieChartIcon, Plus } from 'lucide-react'
import { useFetch } from '@/hooks/util'
import { DateRange, dateRangeMap } from '@/client/trends/type'
import LanguageSelector from './language-selector'
import { Badge } from '@/components/ui/badge'
import TrendRepository from './trending-repo'
import { ITrendRepository } from '@/client/trends/type'

function capitalizeFirstLetter(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1)
}

export default function TrendRepos() {
  const [since, setSince] = useState<DateRange>('daily')
  const [language, setLanguage] = useState<string | null>(null)
  const { isLoading, run, data } = useFetch<ITrendRepository[]>({
    initialRun: false,
    config: {
      url: '/trending',
      method: 'GET',
      params: {
        since,
        language,
      },
    },
  })

  const selectLanguage = (e: string | null) => {
    setLanguage(e)
  }

  useEffect(() => {
    run({ params: { since, language: language ?? '' } })
  }, [language, since])

  return (
    <Card className='flex flex-col h-[320px] w-full max-w-[380px] md:max-w-none overflow-x-hidden'>
      <CardHeader className='items-start pb-0 gap-y-0'>
        <CardTitle className='text-xl flex gap-x-4 justify-between w-full'>
          <div className='flex items-center gap-x-3'>
            <div>Trending Repositories</div>
            <Badge variant='secondary'>from GitHub</Badge>
          </div>
        </CardTitle>
        <CardDescription className='flex flex-col gap-2 md:flex-row lg:flex-col 2xl:flex-row'>
          <LanguageSelector selected={language} onChange={selectLanguage} disabled={isLoading} />
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
          data.map(repo => (
            <TrendRepository key={`trend_repo_${repo.owner_name}_${repo.repo_name}`} {...repo} date_range={since} />
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
