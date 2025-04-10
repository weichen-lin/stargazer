import { Download, Star, FolderPlus } from 'lucide-react';
import { Button } from '@/components/ui/button';

export default function Welcome() {
  return (
    <div className='row-span-1 flex items-center'>
      <div className='w-full px-4 py-2 flex items-center justify-between bg-card rounded-lg shadow-sm'>
        <div className=''>
          <h1 className='text-xl font-bold'>Welcome back, WeiChen!</h1>
          <p className='text-muted-foreground mt-1'>
            Organize and manage your GitHub stars effortlessly
          </p>
        </div>
        <div className='flex flex-wrap gap-6'>
          <Button>
            <div className='flex items-center gap-x-2'>
              <Download className='h-4 w-4' />
              <span>Import GitHub Stars</span>
            </div>
          </Button>
          <Button variant='outline' className='flex items-center gap-2'>
            <div className='flex gap-x-2 items-center'>
              <FolderPlus className='h-4 w-4' />
              <span>Create New Group</span>
            </div>
          </Button>
        </div>
      </div>
    </div>
  );
}
