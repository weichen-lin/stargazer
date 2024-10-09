import { Switch } from '@/components/ui/switch'
import { cn } from '@/lib/utils'

export default function Switcher() {
  return (
    <div
      className={cn(
        'flex items-center space-x-2 transition-colors border-[1px]',
        'bg-white rounded-lg border-slate-300 shadow-md px-3 py-1',
      )}
    >
      <Switch />
      <div className='text-sm'>Make it public</div>
    </div>
  )
}
