'use server'

import { SelectLanguage, Results } from './components'
import { StarsProvider } from './hook/useStarsContext'

export default async function Stars() {
  return (
    <div className='p-3 lg:p-6 w-full flex flex-col h-full overflow-hidden'>
      <StarsProvider>
        <div className='w-full'>
          <SelectLanguage />
        </div>
        <div className='w-full flex flex-col items-center justify-center flex-1 overflow-y-auto mb-8'>
          <Results />
        </div>
      </StarsProvider>
    </div>
  )
}
