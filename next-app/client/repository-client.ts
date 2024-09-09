'use server'

import BaseClient from './base-client'
export interface ILanguageDistribution {
  language: string
  count: number
}

class RepositoryClient extends BaseClient {
  constructor(email: string) {
    super(email)
  }

  getLanguageDistribution() {
    return this.get<ILanguageDistribution[]>('/repository/language-distribution')
  }
}

export default RepositoryClient
