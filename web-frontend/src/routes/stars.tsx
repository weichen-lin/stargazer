import { createFileRoute } from '@tanstack/react-router';
import AuthLayout from '@/components/layout/auth';
import { Search } from '@/pages/stars';

export const Route = createFileRoute('/stars')({
  component: Stars,
});

function Stars() {
  return (
    <AuthLayout>
      <Search />
    </AuthLayout>
  );
}
