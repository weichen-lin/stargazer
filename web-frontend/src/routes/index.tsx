import { createFileRoute } from '@tanstack/react-router';
import AuthLayout from '@/components/layout/auth';

export const Route = createFileRoute('/')({
  component: Index,
});

function Index() {
  return <AuthLayout>Index</AuthLayout>;
}
