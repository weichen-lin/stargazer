import BaseClient from '@/client/base-client'
import { ICollection, ISharedCollection, UpdateCollectionPayload } from './type'

class CollectionClient extends BaseClient {
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

  getReposCollections(id: string, page: string, limit: string) {
    const params = new URLSearchParams({
      page,
      limit,
    })

    return this.get<ICollection[]>(`/collection/repos/${id}`, params)
  }

  createCollection(name: string) {
    return this.post<any, ICollection>('/collection', {
      name,
    })
  }

  deleteCollection(id: string) {
    return this.delete(`/collection`, { id })
  }

  updateCollection(id: string, payload: UpdateCollectionPayload) {
    return this.patch<UpdateCollectionPayload, ICollection>(`/collection/${id}`, payload)
  }

  addRepoToCollection(id: string, repo_ids: number[]) {
    return this.post<any, ICollection>(`/collection/repos/${id}`, {
      repo_ids,
    })
  }
}

export default CollectionClient
