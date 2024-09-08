'use server'

import BaseClient from './base-client'

class RepositoryClient extends BaseClient {
  constructor(email: string) {
    super(email)
  }

  getLanguageDistribution() {
    return this.get('/repository/language-distribution')
  }
}

export default RepositoryClient
