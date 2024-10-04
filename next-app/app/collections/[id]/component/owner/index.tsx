import { ISharedFrom } from '@/client/collection/type'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { MessageCircleHeartIcon } from 'lucide-react'
import { Alert, AlertDescription } from '@/components/ui/alert'

export default function Owner(props: ISharedFrom) {
  const { name, image, email } = props

  return (
    <Alert className='max-w-[380px] flex flex-col shadow-md'>
      <AlertDescription className='flex gap-x-2 items-center w-full'>
        <MessageCircleHeartIcon className='w-6 h-6 text-red-300' />
        <div>Created by</div>
        <div className='flex items-center gap-2'>
          <Avatar>
            <AvatarImage src={image} />
            <AvatarFallback>{name.substring(0, 2).toUpperCase()}</AvatarFallback>
          </Avatar>
          <div className='text-sm'>{name}</div>
        </div>
      </AlertDescription>
    </Alert>
  )
}
