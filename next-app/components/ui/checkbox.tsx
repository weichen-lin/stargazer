import { cn } from '@/lib/utils'
import { useEffect, useState } from 'react'

interface ICheckboxProps {
  checked?: boolean
  onChange?: (checked: boolean) => void
}

export function Checkbox({ checked, onChange }: ICheckboxProps) {
  const [isCheck, setIsChecked] = useState(checked ?? false)

  return (
    <label className='relative block cursor-pointer select-none rounded-full text-2xl outline-2 outline-offset-1 outline-[#385cb0] has-[:checked]:rounded-md has-[:focus]:outline'>
      <input className='peer absolute h-0 w-0 opacity-0' type='checkbox' />
      <div
        className={cn(
          'relative left-0 top-0 h-[1.5rem] w-[1.5rem] rounded-[50%] bg-slate-200 transition duration-300 focus:outline-[#385cb0]',
          'peer-checked:rounded-lg peer-checked:bg-[#385cb0] peer-checked:after:block',
          'after:absolute after:left-[0.5rem] after:top-1 after:hidden after:h-[0.8rem] after:w-[0.5rem] after:rotate-45 after:border-b-[0.2rem]',
          'after:border-r-[0.2rem] after:content-[""]',
        )}
      />
    </label>
  )
}
