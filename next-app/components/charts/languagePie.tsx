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

  return !isLoaded ? (
    <div className='pt-3 pl-3 mx-auto w-[380px] h-[380px] border-[2px] border-slate-300/70 rounded-lg bg-white drop-shadow-md'>
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
    </div>
  ) : (
    <LanguagePieLoading />
  )
}

const LanguagePieLoading = () => {
  return <div>asdasd</div>
}
