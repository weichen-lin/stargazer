import BaseClient from '@/client/base-client'
import { z } from 'zod'
import { ICollection } from './type'

class CollectionClient extends BaseClient {
  constructor(email: string) {
    super(email)
  }

  getCollections(page: string, limit: string) {
    const params = new URLSearchParams({
      page,
      limit,
    })

    return this.get<{ total: number; data: ICollection[] }>('/collection', params)
  }

  createCollection(name: string) {
    return this.post<any, ICollection>('/collection', {
      name,
    })
  }
}

export default CollectionClient
