'use server'

import { HomeLayout } from '@/components/layout'

interface LayoutProps {
  children: React.ReactNode
  params: { username: string }
}

export default async function Layout(props: Readonly<LayoutProps>) {
  const { username } = props.params

  return <HomeLayout path='settings'>{props.children}</HomeLayout>
}
