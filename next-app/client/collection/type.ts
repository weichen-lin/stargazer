export interface ICollection {
  id: string
  name: string
  is_public: boolean
  created_at: string
  updated_at: string
}

export interface ISharedFrom {
  name: string
  image: string
  email: string
}

export interface ISharedCollection {
  owner: string
  collection: ICollection
  shared_from: ISharedFrom | null
}
