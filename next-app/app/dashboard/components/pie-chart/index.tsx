'use client'

import { useMemo, useEffect } from 'react'
import { Label, Pie, PieChart } from 'recharts'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ChartContainer, ChartTooltip, ChartTooltipContent } from '@/components/ui/chart'
import { colorConfig, getLanguageColor, Language } from './config'
import { PieChart as PieChartIcon, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useUser } from '@/context'

const testLanguages = [
  {
    language: 'TypeScript',
    count: 184,
  },
  {
    language: 'Go',
    count: 103,
  },
  {
    language: 'Python',
    count: 98,
  },
  {
    language: 'JavaScript',
    count: 88,
  },
  {
    language: 'Unknown',
    count: 33,
  },
  {
    language: 'Rust',
    count: 21,
  },
  {
    language: 'Shell',
    count: 10,
  },
  {
    language: 'C',
    count: 10,
  },
  {
    language: 'Java',
    count: 10,
  },
  {
    language: 'HTML',
    count: 9,
  },
  {
    language: 'C++',
    count: 7,
  },
  {
    language: 'PHP',
    count: 3,
  },
  {
    language: 'Jupyter Notebook',
    count: 3,
  },
  {
    language: 'Swift',
    count: 3,
  },
  {
    language: 'Clojure',
    count: 2,
  },
  {
    language: 'Markdown',
    count: 2,
  },
  {
    language: 'MDX',
    count: 2,
  },
  {
    language: 'Dockerfile',
    count: 2,
  },
  {
    language: 'Vue',
    count: 2,
  },
  {
    language: 'OCaml',
    count: 1,
  },
  {
    language: 'SCSS',
    count: 1,
  },
  {
    language: 'HCL',
    count: 1,
  },
  {
    language: 'Svelte',
    count: 1,
  },
  {
    language: 'Ruby',
    count: 1,
  },
  {
    language: 'Dart',
    count: 1,
  },
  {
    language: 'SVG',
    count: 1,
  },
  {
    language: 'Lua',
    count: 1,
  },
  {
    language: 'Jinja',
    count: 1,
  },
  {
    language: 'CSS',
    count: 1,
  },
  {
    language: 'Pug',
    count: 1,
  },
]

export default function LanguageDistribution() {
  const { email } = useUser()
  const totalStars = useMemo(() => {
    return testLanguages.reduce((acc, curr) => acc + curr.count, 0)
  }, [])

  const data = testLanguages.map(e => ({
    language: e.language,
    count: e.count,
    fill: getLanguageColor(e.language as Language),
  }))

  useEffect(() => {
    const res = async () => {
      const res = await fetch('/api/repository/language-distribution')
      console.log({ res })
    }

    res()
  }, [])

  return (
    <Card className='flex flex-col h-[320px] w-full max-w-[380px] md:max-w-none'>
      <CardHeader className='items-center pb-0 gap-y-1'>
        <CardTitle className='text-xl'>Language Distribution</CardTitle>
      </CardHeader>
      <CardContent className='flex-1'>
        <ChartContainer className='mx-auto aspect-square max-h-[250px] py-4' config={colorConfig}>
          <PieChart>
            <ChartTooltip cursor={false} content={<ChartTooltipContent hideLabel />} />
            <Pie data={data} dataKey='count' nameKey='language' innerRadius={60} strokeWidth={5}>
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
          {/* <EmptyContent /> */}
        </ChartContainer>
      </CardContent>
      {/* <Loading /> */}
    </Card>
  )
}

const Loading = () => {
  return (
    <div className='flex flex-col items-center justify-center py-4 gap-y-4'>
      <div className='w-[200px] h-[200px] rounded-full bg-slate-200 animate-pulse'></div>
      <div className='w-2/3 rounded-full bg-slate-200 animate-pulse h-8'></div>
    </div>
  )
}

const EmptyContent = () => {
  return (
    <div className='w-full flex flex-col items-center justify-center'>
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
