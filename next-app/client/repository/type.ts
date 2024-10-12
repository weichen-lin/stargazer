import { z } from 'zod'
import { ICollection } from '../collection'

export interface ILanguageDistribution {
  language: string
  count: number
}

export interface ITopics {
  name: string
  repos: number[]
}

export const SortKeySchema = z.enum(['created_at', 'stargazers_count', 'watchers_count'])

export type SortKey = z.infer<typeof SortKeySchema>

export interface IRepository {
  repo_id: number
  repo_name: string
  owner_name: string
  avatar_url: string
  html_url: string
  homepage?: string
  description: string
  created_at: string
  updated_at: string
  stargazers_count: number
  language: string
  watchers_count: number
  open_issues_count: number
  default_branch: string
  archived: boolean
  topics: string[]
  external_created_at: string
  last_synced_at: string
  last_modified_at: string
}

export interface IRepoSearchWithLanguage {
  total: number
  data: { repository: IRepository; collected_by: ICollection[] }[]
}
