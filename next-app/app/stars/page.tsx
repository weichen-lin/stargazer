'use server'

import { SelectLanguage } from './components'

export default async function Stars() {
  return (
    <div className='py-6'>
      <div className='max-w-[400px]'>
        <SelectLanguage />
      </div>
    </div>
  )
}
