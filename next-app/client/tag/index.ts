import BaseClient from '@/client/base-client'
import { ITag } from './type'

export interface ITagPayload {
  name: string
  repo_id: number
}

class TagClient extends BaseClient {
  getTags(id: string) {
    return this.get<ITag[]>(`/tag/${id}`)
  }

  createTag(req: ITagPayload) {
    return this.post<ITagPayload, any>('/tag', req)
  }

  deleteTag(req: ITagPayload) {
    return this.delete<ITagPayload, any>('/tag', req)
  }
}

export type { ITag }
export default TagClient
