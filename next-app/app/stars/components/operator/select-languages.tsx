'use client'

import MultipleSelector from '@/components/ui/multiple-selector'
import { Button } from '@/components/ui/button'
import { FixPagination } from '@/components/tab'
import { ILanguageDistribution } from '@/client/repository'
import { useFetch } from '@/hooks/util'
import SelectRepository from './select-repository'
import { useStarsContext } from '@/app/stars/hook'

const SelectLanguage = () => {
  const { data, isLoading } = useFetch<ILanguageDistribution[]>({
    initialRun: true,
    config: {
      url: '/repository/language-distribution',
      method: 'GET',
    },
  })

  const { selectLanguages, setSelectLanguages, search, total } = useStarsContext()

  const options = (data && data.map(e => ({ label: e.language, value: e.language }))) || []

  return (
    <div className='flex flex-col gap-2 items-center lg:items-start w-[95%] lg:w-full lg:px-4 mx-auto'>
      <div className='flex flex-col gap-y-2 lg:w-[380px]'>
        <div className=''>Language</div>
        {isLoading ? (
          <div className='w-[380px] h-10 bg-slate-300 rounded-lg animate-pulse'></div>
        ) : (
          <MultipleSelector
            value={selectLanguages}
            onChange={setSelectLanguages}
            defaultOptions={options}
            placeholder='Select languages you like...'
            emptyIndicator={
              <p className='text-center text-lg leading-10 text-gray-600 dark:text-gray-400'>no results found.</p>
            }
            className='border-slate-900 w-[380px]'
            disabled={isLoading}
          />
        )}
      </div>
      <div className='flex justify-between w-[380px]'>
        <Button
          onClick={() => {
            search(1)
          }}
          className='w-20 h-10 mb-2'
          disabled={isLoading || selectLanguages.length === 0}
        >
          Search
        </Button>
        <SelectRepository />
      </div>
      {total > 0 && <FixPagination total={total} />}
    </div>
  )
}

export default SelectLanguage
