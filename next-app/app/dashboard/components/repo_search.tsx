'use client'

import clsx from 'clsx'
import { getReposByKey, ISearchKey } from '@/actions/neo4j'
import Repo from './repo'
import { useEffect, useState } from 'react'
import { IRepoAtDashboard } from '@/actions/neo4j'
import { useUser } from '@/context'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'

const fakeData = [
  {
    owner_name: 'ClickHouse',
    stargazers_count: 36484,
    open_issues_count: 3896,
    created_at: '2016-06-02T08:28:18Z',
    description:
      'ClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMS',
    language: 'Go',
    archived: false,
    avatar_url: 'https://avatars.githubusercontent.com/u/54801242?v=4',
    updated_at: '2024-09-01T10:12:20Z',
    repo_id: 60246359,
    html_url: 'https://github.com/ClickHouse/ClickHouse',
    default_branch: 'master',
    repo_name: 'ClickHouse',
    watchers_count: 36484,
    homepage: 'https://clickhouse.com',
  },
  {
    owner_name: 'ClickHouse',
    stargazers_count: 36484,
    open_issues_count: 3896,
    created_at: '2016-06-02T08:28:18Z',
    description:
      'ClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMS',
    language: 'Go',
    archived: false,
    avatar_url: 'https://avatars.githubusercontent.com/u/54801242?v=4',
    updated_at: '2024-09-01T10:12:20Z',
    repo_id: 60246359,
    html_url: 'https://github.com/ClickHouse/ClickHouse',
    default_branch: 'master',
    repo_name: 'ClickHouse',
    watchers_count: 36484,
    homepage: 'https://clickhouse.com',
  },
  {
    owner_name: 'ClickHouse',
    stargazers_count: 36484,
    open_issues_count: 3896,
    created_at: '2016-06-02T08:28:18Z',
    description:
      'ClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMS',
    language: 'Go',
    archived: false,
    avatar_url: 'https://avatars.githubusercontent.com/u/54801242?v=4',
    updated_at: '2024-09-01T10:12:20Z',
    repo_id: 60246359,
    html_url: 'https://github.com/ClickHouse/ClickHouse',
    default_branch: 'master',
    repo_name: 'ClickHouse',
    watchers_count: 36484,
    homepage: 'https://clickhouse.com',
  },
  {
    owner_name: 'ClickHouse',
    stargazers_count: 36484,
    open_issues_count: 3896,
    created_at: '2016-06-02T08:28:18Z',
    description:
      'ClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMS',
    language: 'Go',
    archived: false,
    avatar_url: 'https://avatars.githubusercontent.com/u/54801242?v=4',
    updated_at: '2024-09-01T10:12:20Z',
    repo_id: 60246359,
    html_url: 'https://github.com/ClickHouse/ClickHouse',
    default_branch: 'master',
    repo_name: 'ClickHouse',
    watchers_count: 36484,
    homepage: 'https://clickhouse.com',
  },
  {
    owner_name: 'ClickHouse',
    stargazers_count: 36484,
    open_issues_count: 3896,
    created_at: '2016-06-02T08:28:18Z',
    description:
      'ClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMS',
    language: 'Go',
    archived: false,
    avatar_url: 'https://avatars.githubusercontent.com/u/54801242?v=4',
    updated_at: '2024-09-01T10:12:20Z',
    repo_id: 60246359,
    html_url: 'https://github.com/ClickHouse/ClickHouse',
    default_branch: 'master',
    repo_name: 'ClickHouse',
    watchers_count: 36484,
    homepage: 'https://clickhouse.com',
  },
  {
    owner_name: 'ClickHouse',
    stargazers_count: 36484,
    open_issues_count: 3896,
    created_at: '2016-06-02T08:28:18Z',
    description:
      'ClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMS',
    language: 'Go',
    archived: false,
    avatar_url: 'https://avatars.githubusercontent.com/u/54801242?v=4',
    updated_at: '2024-09-01T10:12:20Z',
    repo_id: 60246359,
    html_url: 'https://github.com/ClickHouse/ClickHouse',
    default_branch: 'master',
    repo_name: 'ClickHouse',
    watchers_count: 36484,
    homepage: 'https://clickhouse.com',
  },
  {
    owner_name: 'ClickHouse',
    stargazers_count: 36484,
    open_issues_count: 3896,
    created_at: '2016-06-02T08:28:18Z',
    description:
      'ClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMS',
    language: 'Go',
    archived: false,
    avatar_url: 'https://avatars.githubusercontent.com/u/54801242?v=4',
    updated_at: '2024-09-01T10:12:20Z',
    repo_id: 60246359,
    html_url: 'https://github.com/ClickHouse/ClickHouse',
    default_branch: 'master',
    repo_name: 'ClickHouse',
    watchers_count: 36484,
    homepage: 'https://clickhouse.com',
  },
  {
    owner_name: 'ClickHouse',
    stargazers_count: 36484,
    open_issues_count: 3896,
    created_at: '2016-06-02T08:28:18Z',
    description:
      'ClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMSClickHouse® is a real-time analytics DBMS',
    language: 'Go',
    archived: false,
    avatar_url: 'https://avatars.githubusercontent.com/u/54801242?v=4',
    updated_at: '2024-09-01T10:12:20Z',
    repo_id: 60246359,
    html_url: 'https://github.com/ClickHouse/ClickHouse',
    default_branch: 'master',
    repo_name: 'ClickHouse',
    watchers_count: 36484,
    homepage: 'https://clickhouse.com',
  },
]

export default function RepoSearch({ searchKey, title }: { searchKey: ISearchKey; title: string }) {
  const [repos, setRepos] = useState<IRepoAtDashboard[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const { email } = useUser()

  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await getReposByKey(email, searchKey)
        setRepos(res)
      } catch (error) {
        console.error(error)
      } finally {
        setIsLoading(false)
      }
    }
    fetchData()
  }, [])

  return (
    <Card className='flex flex-col h-[320px]'>
      <CardHeader className='items-start pb-0 gap-y-0'>
        <CardTitle className='text-xl'>Recent Activity</CardTitle>
        <CardDescription>Your latest starred repositories</CardDescription>
      </CardHeader>
      <CardContent className='flex-1 pb-0 overflow-y-scroll flex flex-col gap-y-2 py-4'>
        {fakeData.map(e => (
          <Repo {...e} searchKey='' key={`${e.repo_id}`} />
        ))}
      </CardContent>
    </Card>
  )
}
