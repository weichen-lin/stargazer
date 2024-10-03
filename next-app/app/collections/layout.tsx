'use server'

import { AuthLayout } from '@/components/layout'

interface LayoutProps {
  children: React.ReactNode
}

export default async function Layout({ children }: Readonly<LayoutProps>) {
  return <AuthLayout>{children}</AuthLayout>
}
