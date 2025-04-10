import { createFileRoute } from '@tanstack/react-router';
import AuthLayout from '@/components/layout/auth';
import { Download, Star, FolderPlus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { Welcome, MyCollections, RecentStars } from '@/pages/dashboard';

export const Route = createFileRoute('/dashboard')({
  component: Dashboard,
});

function Dashboard() {
  return (
    <AuthLayout>
      <div className='grid grid-rows-9 grid-flow-col gap-4 h-full px-4 py-2'>
        <Welcome />
        <MyCollections />
        <RecentStars />
      </div>
    </AuthLayout>
  );
}
