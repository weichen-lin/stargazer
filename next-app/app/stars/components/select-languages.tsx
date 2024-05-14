'use client'

import MultipleSelector, { Option } from '@/components/ui/multiple-selector'
import { getLanguageDistribution, ILanguageDistribution } from '@/actions/neo4j'
import { useState, useEffect } from 'react'
import { useUser } from '@/context'
import { useStars } from '@/hooks/stars'
import { Button } from '@/components/ui/button'
import { FixPagination } from '@/components/tab'

const SelectLanguage = () => {
  const { name } = useUser()
  const [isLoaded, setIsLoaded] = useState(true)
  const [data, setData] = useState<Option[]>([])
  const { selected, setSelected, search, count } = useStars()

  useEffect(() => {
    const fetchData = async () => {
      const res = await getLanguageDistribution(name)
      const data = res.map((item: ILanguageDistribution) => {
        if (item.language === '') {
          return { label: 'unknown', value: '' }
        }
        return { label: item.language, value: item.language }
      })

      setData(data)
      setIsLoaded(false)
    }
    fetchData()
  }, [])

  return (
    <div className='flex gap-x-8 items-end'>
      <div className='flex flex-col gap-y-4'>
        <div>Language</div>
        {isLoaded ? (
          <div className='w-[400px] h-10 bg-slate-300 rounded-lg animate-pulse'></div>
        ) : (
          <MultipleSelector
            value={selected}
            onChange={setSelected}
            defaultOptions={data}
            placeholder='Select languages you like...'
            emptyIndicator={
              <p className='text-center text-lg leading-10 text-gray-600 dark:text-gray-400'>no results found.</p>
            }
            className='border-slate-900 w-[400px]'
            disabled={isLoaded}
          />
        )}
      </div>
      <Button
        onClick={() => {
          search(1)
        }}
        className='w-20 h-10 mb-2'
      >
        Search
      </Button>
      {count > 0 && <FixPagination total={count} />}
    </div>
  )
}

export default SelectLanguage
