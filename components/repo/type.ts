export interface IStar {
  id: number
  full_name: string
  owner: {
    avatar_url: string
  }
  html_url: string
  description: string
  homepage: string
  stargazers_count: number
  language: string
}
