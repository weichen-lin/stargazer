'use client'

import { AgChartsReact } from 'ag-charts-react'
import { useUser } from '@/context'
import { useState, useEffect } from 'react'
import { getLanguageDistribution, ILanguageDistribution } from '@/actions/neo4j'
import { useMedia } from '@/hooks/util'
import clsx from 'clsx'

export default function LanguagePie() {
  const { name } = useUser()
  const [isLoaded, setIsLoaded] = useState(true)
  const [data, setData] = useState<ILanguageDistribution[]>([])

  const { isRetina } = useMedia()

  useEffect(() => {
    const fetchData = async () => {
      const res = await getLanguageDistribution(name)
      setData(res)
      setIsLoaded(false)
    }
    fetchData()
  }, [])

  return (
    <div
      className={clsx(
        'pt-3 pl-3 xl:pt-1 xl:pl-1 2xl:pt-3 2xl:pl-3 border-[2px] border-slate-300/70 rounded-lg bg-white drop-shadow-md',
        'w-[380px] h-[380px] xl:w-[330px] 2xl:w-[380px] xl:h-[330px] 2xl:h-[380px]',
      )}
    >
      {isLoaded ? (
        <LanguagePieLoading />
      ) : data.length > 0 ? (
        <AgChartsReact
          options={{
            data: data,
            title: {
              text: 'Liked Languages Distribution',
            },
            series: [
              {
                type: 'pie',
                angleKey: 'count',
                calloutLabelKey: 'language',
                sectorLabelKey: 'language',
                sectorLabel: {
                  color: 'white',
                  fontWeight: 'bold',
                },
              },
            ],
            width: isRetina ? 320 : 350,
            height: isRetina ? 320 : 350,
          }}
        />
      ) : (
        <>No current data.</>
      )}
    </div>
  )
}

const LanguagePieLoading = () => {
  return (
    <div className='flex flex-col w-full h-full items-center justify-center gap-y-12'>
      <div className='w-1/5 h-6 bg-slate-300 rounded-md animate-pulse'></div>
      <div className='w-3/5 h-12 bg-slate-300 rounded-md animate-pulse'></div>
    </div>
  )
}
