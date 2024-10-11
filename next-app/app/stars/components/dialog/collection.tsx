import { ICollection } from '@/client/collection'
import { MixIcon } from '@radix-ui/react-icons'

export default function ListCollection() {
  //   const { id, name } = props

  return (
    <div className='flex justify-between h-8 hover:bg-slate-100 w-full px-4 cursor-pointer py-2'>
      <div className='flex gap-x-2 items-center'>
        <MixIcon className='w-5 h-5' />
        <div>test collectikon</div>
      </div>
      <div></div>
    </div>
  )
}
