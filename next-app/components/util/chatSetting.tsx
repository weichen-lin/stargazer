'use client'

import { InfoCircledIcon, GearIcon } from '@radix-ui/react-icons'
import { Button } from '@/components/ui/button'
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/components/ui/drawer'
import { Dialog, DialogContent, DialogTrigger } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Slider } from '@/components/ui/slider'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { useConfig } from '@/hooks/setting'

export function ChatSettingDialog() {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant='outline' size='icon'>
          <GearIcon className='h-5 w-5' />
        </Button>
      </DialogTrigger>
      <DialogContent className='sm:max-w-[425px] max-w-[375px] mx-auto'>
        <ChatSetting />
      </DialogContent>
    </Dialog>
  )
}

export function ChatSettingDrawer() {
  return (
    <Drawer>
      <DrawerTrigger asChild>
        <Button variant='outline' size='icon'>
          <GearIcon className='w-5 h-5' />
          <span className='sr-only'>chat settings</span>
        </Button>
      </DrawerTrigger>
      <DrawerContent>
        <ChatSetting />
      </DrawerContent>
    </Drawer>
  )
}

const ChatSetting = () => {
  const { isLoading, config, change, update } = useConfig()

  return (
    <Tabs defaultValue='tokens' className='flex flex-col items-center w-full pt-6'>
      <TabsList>
        <TabsTrigger value='tokens' disabled={isLoading}>
          Token
        </TabsTrigger>
        <TabsTrigger value='preference' disabled={isLoading}>
          Preference
        </TabsTrigger>
      </TabsList>
      <TabsContent value='tokens'>
        <div className='mx-auto w-[350px] h-[350px] flex flex-col justify-between pb-10'>
          <DrawerHeader>
            <DrawerTitle className='text-2xl text-center'>Setup all token you need</DrawerTitle>
          </DrawerHeader>
          <div className='p-4 w-full flex flex-col gap-y-8'>
            <div className='flex flex-col gap-y-2'>
              <h1 className='font-semibold flex gap-x-2 items-center'>
                Github Token
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger>
                      <InfoCircledIcon color='#CE2C31' />
                    </TooltipTrigger>
                    <TooltipContent>
                      Your GitHub token will only be used to fetch repositories you follow on GitHub.
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </h1>
              <Input
                placeholder='Your Github Token'
                value={config.githubToken}
                onChange={e => {
                  change('githubToken', e.target.value)
                }}
                disabled={isLoading}
              />
            </div>
          </div>
          <div className='p-4 w-full flex flex-col gap-y-8'>
            <div className='flex flex-col gap-y-2'>
              <h1 className='font-semibold flex gap-x-2 items-center'>
                OpenAI API Key
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger>
                      <InfoCircledIcon color='#CE2C31' />
                    </TooltipTrigger>
                    <TooltipContent>
                      Your OpenAI key will only be used for vector search and processing the repositories you follow on
                      GitHub into vectors.
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </h1>
              <Input
                placeholder='Your OpenAI API Key'
                value={config.openAIKey}
                onChange={e => {
                  change('openAIKey', e.target.value)
                }}
                disabled={isLoading}
              />
            </div>
          </div>
          <DrawerFooter>
            <Button
              onClick={async () => {
                await update()
              }}
              loading={isLoading}
              disabled={isLoading}
            >
              Save
            </Button>
            <DrawerClose asChild>
              <Button variant='outline'>Cancel</Button>
            </DrawerClose>
          </DrawerFooter>
        </div>
      </TabsContent>
      <TabsContent value='preference'>
        <div className='mx-auto w-[350px] h-[350px] flex flex-col justify-between pb-10'>
          <DrawerHeader>
            <DrawerTitle className='text-2xl text-center'>Search Preference</DrawerTitle>
          </DrawerHeader>
          <div className='p-4 w-full flex flex-col gap-y-8'>
            <div className='flex flex-col gap-y-4 mb-6'>
              <h1 className='font-semibold flex items-center gap-x-4'>
                Cosine Similarity : <span>{config.cosine}</span>
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger>
                      <InfoCircledIcon color='#CE2C31' />
                    </TooltipTrigger>
                    <TooltipContent>Will only be adopted if it exceeds this value.</TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </h1>
              <Slider
                defaultValue={[config.cosine]}
                max={1}
                step={0.01}
                onValueChange={e => {
                  change('cosine', e[0])
                }}
                disabled={isLoading}
              />
            </div>
            <div className='flex flex-col gap-y-2'>
              <h1 className='font-semibold'>Response Limit</h1>
              <Input
                placeholder='the limit of response is 20'
                type='number'
                value={config.limit}
                onChange={e => {
                  change('limit', parseInt(e.target.value))
                }}
                onBlur={() => {
                  if (config.limit > 20) {
                    change('limit', 20)
                  }
                }}
                max={20}
                disabled={isLoading}
              />
            </div>
          </div>
          <DrawerFooter>
            <Button
              onClick={async () => {
                await update()
              }}
              loading={isLoading}
              disabled={isLoading}
            >
              Save
            </Button>
            <DrawerClose asChild>
              <Button variant='outline'>Cancel</Button>
            </DrawerClose>
          </DrawerFooter>
        </div>
      </TabsContent>
    </Tabs>
  )
}
