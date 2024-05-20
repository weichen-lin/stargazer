import { GitHubLogoIcon, LightningBoltIcon } from '@radix-ui/react-icons'
import { IRepoAtDashboard } from '@/actions/neo4j'
import moment from 'moment'

export default function Repo(props: IRepoAtDashboard & { searchKey: string }) {
  const { repo_id, html_url, full_name, avatar_url, created_at, last_updated_at, open_issues_count, searchKey } = props

  const date =
    searchKey === 'created_at' ? moment(created_at).format('YYYY/MM/DD') : moment(last_updated_at).format('YYYY/MM/DD')

  return (
    <div className='w-full flex justify-between py-2'>
      <div className='flex gap-x-2 items-center'>
        <img width={32} height={32} src={avatar_url} alt={full_name} />
        <div className='text-slate-700 max-w-[200px] xl:max-w-[140px] 2xl:max-w-[200px] truncate'>{full_name}</div>
      </div>
      <div className='flex gap-x-2 items-center justify-center'>
        {searchKey === 'open_issues_count' ? (
          <div className='flex gap-x-1 items-center justify-center'>
            <LightningBoltIcon className='text-yellow-500' />
            <div className='text-sm pt-1 xl:pt-[2px] 2xl:pt-1'>{open_issues_count}</div>
          </div>
        ) : (
          <div className='text-sm pt-1 xl:pt-[2px] 2xl:pt-1'>{date}</div>
        )}
        <a className='rounded-full hover:bg-slate-300/30' href={html_url} target='_blank'>
          <GitHubLogoIcon />
        </a>
      </div>
    </div>
  )
}
