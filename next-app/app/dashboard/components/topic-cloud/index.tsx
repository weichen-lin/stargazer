import { Text } from '@visx/text'
import { scaleLog } from '@visx/scale'
import Wordcloud from '@visx/wordcloud/lib/Wordcloud'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { generateRandomColors } from './helper'
import { useFetch } from '@/hooks/util'
import { ITopics } from '@/client/repository'
import { PieChart as PieChartIcon, Plus } from 'lucide-react'

export interface WordData {
  text: string
  value: number
}

const colors = generateRandomColors()

function getRotationDegree() {
  const rand = Math.random()
  const degree = rand > 0.5 ? 60 : -60
  return rand * degree
}

const fixedValueGenerator = () => 0.5

export default function TopicsCloud() {
  const { data, isLoading } = useFetch<{ data: ITopics[] }>({
    initialRun: true,
    config: {
      url: '/repository/topics',
      method: 'GET',
    },
  })

  const word_data = data?.data ?? []

  const words =
    word_data
      ?.sort((a, b) => b.repos.length - a.repos.length)
      .slice(0, 50)
      .map(e => ({
        text: e.name,
        value: e.repos.length,
      })) ?? []

  const fontScale = scaleLog({
    domain: [Math.min(...words.map(w => w.value)), Math.max(...words.map(w => w.value))],
    range: [10, 100],
  })

  const fontSizeSetter = (datum: WordData) => fontScale(datum.value)

  return (
    <div className='flex flex-col h-[320px] w-full max-w-[380px] md:max-w-none'>
      <Card className='flex flex-col h-[320px] w-full'>
        <CardHeader className='items-center pb-0 gap-y-1'>
          <CardTitle className='text-xl'>Topics Cloud</CardTitle>
        </CardHeader>
        <CardContent className='flex-1 flex items-center justify-center'>
          {!isLoading && word_data && word_data.length > 0 && (
            <Wordcloud
              words={words}
              width={340}
              height={260}
              fontSize={fontSizeSetter}
              font={'Impact'}
              padding={2}
              spiral='rectangular'
              rotate={getRotationDegree}
              random={fixedValueGenerator}
            >
              {cloudWords =>
                cloudWords.map((w, i) => (
                  <Text
                    key={w.text}
                    fill={colors[i % colors.length]}
                    textAnchor={'middle'}
                    transform={`translate(${w.x}, ${w.y}) rotate(${w.rotate})`}
                    fontSize={w.size}
                    fontFamily={w.font}
                    onMouseEnter={e => {}}
                  >
                    {w.text}
                  </Text>
                ))
              }
            </Wordcloud>
          )}
          {!isLoading && word_data && word_data.length === 0 && <EmptyContent />}
        </CardContent>
        {isLoading && <Loading />}
      </Card>
    </div>
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
    <div className='w-full flex flex-col items-center justify-center my-1'>
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
