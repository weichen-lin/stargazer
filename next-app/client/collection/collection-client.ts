import BaseClient from '@/client/base-client'
import { ICollection, ISharedCollection } from './type'

class CollectionClient extends BaseClient {
  constructor(email: string) {
    super(email)
  }

  getCollection(id: string) {
    return this.get<ISharedCollection>(`/collection/${id}`)
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

  deleteCollection(id: string) {
    return this.delete(`/collection`, { id })
  }
}

export default CollectionClient
