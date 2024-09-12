'use client'

import { useMemo } from 'react'
import { Label, Pie, PieChart } from 'recharts'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ChartContainer, ChartTooltip, ChartTooltipContent } from '@/components/ui/chart'
import { colorConfig, getLanguageColor, Language } from './config'
import { PieChart as PieChartIcon, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useFetch } from '@/hooks/util'
import { ILanguageDistribution } from '@/client/repository'

export default function LanguageDistribution() {
  const { data, isLoading } = useFetch<ILanguageDistribution[]>({
    initialRun: true,
    config: {
      url: '/repository/language-distribution',
      method: 'GET',
    },
  })

  const totalStars = useMemo(() => {
    return data ? data.reduce((acc, curr) => acc + curr.count, 0) : 0
  }, [data])

  const withColor = data
    ? data.map(e => ({
        language: e.language,
        count: e.count,
        fill: getLanguageColor(e.language as Language),
      }))
    : []

  return (
    <Card className='flex flex-col h-[320px] w-full max-w-[380px] md:max-w-none'>
      <CardHeader className='items-center pb-0 gap-y-1'>
        <CardTitle className='text-xl'>Language Distribution</CardTitle>
      </CardHeader>
      <CardContent className='flex-1'>
        {!isLoading && data && data.length > 0 && (
          <ChartContainer className='mx-auto aspect-square max-h-[250px] py-4' config={colorConfig}>
            <PieChart>
              <ChartTooltip cursor={false} content={<ChartTooltipContent hideLabel />} />
              <Pie data={withColor} dataKey='count' nameKey='language' innerRadius={60} strokeWidth={5}>
                <Label
                  content={({ viewBox }) => {
                    if (viewBox && 'cx' in viewBox && 'cy' in viewBox) {
                      return (
                        <text x={viewBox.cx} y={viewBox.cy} textAnchor='middle' dominantBaseline='middle'>
                          <tspan x={viewBox.cx} y={viewBox.cy} className='fill-foreground text-3xl font-bold'>
                            {totalStars.toLocaleString()}
                          </tspan>
                          <tspan x={viewBox.cx} y={(viewBox.cy || 0) + 24} className='fill-muted-foreground'>
                            Stars
                          </tspan>
                        </text>
                      )
                    }
                  }}
                />
              </Pie>
            </PieChart>
          </ChartContainer>
        )}
        {!isLoading && data && data.length === 0 && <EmptyContent />}
      </CardContent>
      {isLoading && <Loading />}
    </Card>
  )
}

const Loading = () => {
  return (
    <div className='flex flex-col items-center justify-center py-1 gap-y-4 pb-16'>
      <div className='w-[150px] h-[150px] rounded-full bg-slate-200 animate-pulse'></div>
    </div>
  )
}

const EmptyContent = () => {
  return (
    <div className='w-full flex flex-col items-center justify-center my-4'>
      <div className='w-32 h-32 relative'>
        <PieChartIcon className='w-full h-full text-gray-200' />
        <div className='absolute inset-0 flex items-center justify-center'>
          <Plus className='w-8 h-8 text-gray-400' />
        </div>
      </div>
      <p className='text-center text-gray-500 mb-4'>No data yet.</p>
      <Button>
        Start Sync
        <Plus className='w-4 h-4 ml-2' />
      </Button>
    </div>
  )
}
