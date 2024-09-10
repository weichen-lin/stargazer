'use server'

import BaseClient from './base-client'
export interface ILanguageDistribution {
  language: string
  count: number
}

export interface ITopics {
  name: string
  repos: number[]
}

class RepositoryClient extends BaseClient {
  constructor(email: string) {
    super(email)
  }

  getLanguageDistribution() {
    return this.get<ILanguageDistribution[]>('/repository/language-distribution')
  }

  getTopics() {
    return this.get<{ data: ITopics[] }>('/repository/topics')
  }
}

export default RepositoryClient
