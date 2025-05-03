import { Download, FolderPlus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useUser } from '@clerk/clerk-react';
import { toast } from 'sonner';
import useFetch from '@/hooks/useFetch';

export default function Welcome() {
  const { user } = useUser();
  const { run, isLoading } = useFetch({
    config: {
      url: '/background/sync-user-repositories',
      method: 'GET',
    },
    onSuccess: () => {
      toast.success('Successfully imported your GitHub stars!');
    },
    onError: ({ code, message }) => {
      toast.error(message);
    },
  });

  return (
    <div className='row-span-1 flex items-center'>
      <div className='w-full px-4 py-2 flex items-center justify-between bg-card rounded-lg shadow-sm'>
        <div className=''>
          <h1 className='text-xl font-bold'>Welcome back, {user?.username}!</h1>
          <p className='text-muted-foreground mt-1'>
            Organize and manage your GitHub stars effortlessly
          </p>
        </div>
        <Button
          onClick={() => run()}
          variant='outline'
          loading={isLoading}
          className='w-48'
        >
          <div className='flex items-center gap-x-2'>
            <Download className='h-4 w-4' />
            <span>Import GitHub Stars</span>
          </div>
        </Button>
      </div>
    </div>
  );
}
