'use client'

import { AgChartsReact } from 'ag-charts-react'
import { useUser } from '@/context'
import { useState, useEffect } from 'react'
import { getLanguageDistribution, ILanguageDistribution } from '@/actions/neo4j'

export default function LanguagePie() {
  const { name } = useUser()
  const [isLoaded, setIsLoaded] = useState(true)
  const [data, setData] = useState<ILanguageDistribution[]>([])

  useEffect(() => {
    const fetchData = async () => {
      const res = await getLanguageDistribution(name)
      setData(res)
      setIsLoaded(false)
    }
    fetchData()
  }, [])

  return (
    <div className='pt-3 pl-3 w-[380px] h-[380px] border-[2px] border-slate-300/70 rounded-lg bg-white drop-shadow-md'>
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
            width: 350,
            height: 350,
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
