import clsx from 'clsx'
import { TotalRepositoriesLoading, ReposLoading } from './components'

export default async function Loading({ type, title }: { type: string; title: string }) {
  switch (type) {
    case 'total':
      return <TotalRepositoriesLoading />
    default:
      return <ReposLoading title={title} />
  }
}
