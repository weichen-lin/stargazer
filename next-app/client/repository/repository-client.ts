import BaseClient from '@/client/base-client'
import { z } from 'zod'
import { IRepository } from './type'

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

class RepositoryClient extends BaseClient {
  constructor(email: string) {
    super(email)
  }

  syncRepository() {
    return this.get('/repository/sync-repository')
  }

  getRepoDetail(id: string) {
    return this.get<IRepository>(`/repository/${id}`)
  }

  deleteRepo(id: string) {
    return this.delete(`/repository/${id}`)
  }

  getLanguageDistribution() {
    return this.get<ILanguageDistribution[]>('/repository/language-distribution')
  }

  getTopics() {
    return this.get<{ data: ITopics[] }>('/repository/topics')
  }

  getReposWithSortKey(sortkey: SortKey) {
    return this.get<{ data: any[] }>(`/repository/sort?key=${sortkey}&order=DESC`)
  }

  fullTextSearch(query: string) {
    return this.get<IRepository[]>(`/repository/full-text-search?query=${query}`)
  }

  getReposWithLanguage(languages: string, page: string, limit: string) {
    const params = new URLSearchParams({
      languages,
      page,
      limit,
    })

    return this.get<{ total: number; data: IRepository[] }>('/repository', params)
  }
}

export default RepositoryClient
