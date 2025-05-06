import { createFileRoute } from '@tanstack/react-router';
import { Welcome, MyCollections, RecentStars } from '@/pages/dashboard';

export const Route = createFileRoute('/_auth/dashboard')({
  component: Dashboard,
});

function Dashboard() {
  return (
    <div className='grid grid-rows-9 grid-flow-col gap-4 h-full px-4 py-2'>
      <Welcome />
      <MyCollections />
      <RecentStars />
    </div>
  );
}
