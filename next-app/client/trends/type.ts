export type DateRange = 'daily' | 'weekly' | 'monthly'

export const dateRangeMap: { [key in DateRange]: string } = {
  daily: 'today',
  weekly: 'this week',
  monthly: 'this month',
}

export interface ITrendRepository {
  owner_name: string
  repo_name: string
  html_url: string
  description: string
  language: string
  stargazers_count: number
  get_stars: number
}

export interface ITrendDeveloper {
  avatar_url: string
  name: string
  sub_name: string
  html_url: string
  repo_name: string
  description: string
}
