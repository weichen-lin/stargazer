import { TotalLoading, ReposLoading } from './components'

export default async function Loading({ type }: { type: string }) {
  switch (type) {
    case 'total':
      return <TotalLoading />
    default:
      return <ReposLoading />
  }
}
