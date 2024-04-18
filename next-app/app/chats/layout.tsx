'use server'

import { HomeLayout } from '@/components/layout'

interface LayoutProps {
  children: React.ReactNode
}

export default async function Layout(props: Readonly<LayoutProps>) {
  return <HomeLayout path='chats'>{props.children}</HomeLayout>
}
