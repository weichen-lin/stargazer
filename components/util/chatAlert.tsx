'use client'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Cross1Icon, InfoCircledIcon, TriangleRightIcon, GearIcon } from '@radix-ui/react-icons'
import { useChatAlert } from '@/hooks/chat'
import { Button } from '@/components/ui/button'

export default function ChatAlert() {
  const { close } = useChatAlert()
  return (
    <Alert className='lg:w-2/3 lg:mx-auto drop-shadow-md dark:border-slate-100'>
      <div className='flex justify-between items-center'>
        <AlertTitle className='text-lg'>Welcome to StarGazer!</AlertTitle>
        <Cross1Icon className='w-4 h-4 cursor-pointer' onClick={close} />
      </div>
      <AlertDescription>
        To utilize our service, you'll need to provide your own OpenAI API key. This key will solely be used for
      </AlertDescription>
      <AlertDescription className='flex gap-x-2 items-center my-1'>
        <TriangleRightIcon className='w-4 h-4' />
        Vectorizing your starred repositories
      </AlertDescription>
      <AlertDescription className='flex gap-x-2 items-center my-1'>
        <TriangleRightIcon className='w-4 h-4' />
        Search Purpose
      </AlertDescription>
      <AlertDescription className='my-1'>
        Please click on the symbol in the upper right corner
        <GearIcon className='w-5 h-5 inline-block mx-1 mb-1' />
        to provide your API Key.
      </AlertDescription>
      <AlertDescription className='text-gray-700 underline my-1 dark:text-gray-100'>
        The search results include all the data provided by users.
      </AlertDescription>
      <AlertDescription className='flex gap-x-2 items-center my-3'>
        <InfoCircledIcon className='w-4 h-4 mb-[2px]' />
        Your current progress in vectorization
      </AlertDescription>
      <AlertDescription className='flex gap-x-2 items-center my-3'>
        <div className='flex gap-x-2 mx-12'>
          <span>80</span>
          <span>/</span>
          <span>467</span>
        </div>
        <Button onClick={close}>Vectorize</Button>
      </AlertDescription>
    </Alert>
  )
}
