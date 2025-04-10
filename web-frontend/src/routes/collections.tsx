import { createFileRoute } from '@tanstack/react-router';
import AuthLayout from '@/components/layout/auth';

export const Route = createFileRoute('/collections')({
  component: Collections,
});

function Collections() {
  return <AuthLayout>Collections</AuthLayout>;
}
