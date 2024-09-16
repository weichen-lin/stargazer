'use server'

import { load } from 'cheerio'
import axios, { AxiosInstance } from 'axios'

export type DateRange = 'daily' | 'weekly' | 'monthly'

export const dateRangeeMap: { [key in DateRange]: string } = {
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

export const sinceMap: { [key in DateRange]: string } = {
  daily: 'today',
  weekly: 'this week',
  monthly: 'this month',
}

class TrendsClient {
  private readonly _axios!: AxiosInstance
  private readonly _baseUrl: string = 'https://github.com'

  constructor() {
    this._axios = axios.create({
      baseURL: this._baseUrl,
      headers: {
        accept: 'text/html, application/xhtml+xml',
      },
    })
  }

  async getTrendRepos(): Promise<ITrendRepository[]> {
    const res = await this._axios.get('/trending')

    const html = await res.data

    const loader = load(html)

    const nodes = loader('div[data-hpc] article')

    const trendRepos: ITrendRepository[] = []

    for (let i = 0; i < nodes.length; i++) {
      const ch = load(nodes[i])
      const html_url = ch('h2.lh-condensed a').attr('href') ?? ''

      const [owner_name, repo_name] = ch('.lh-condensed a').text().replace(/\s+/g, '').trim().split('/')

      const description = ch('p.color-fg-muted').text().replace(/\n+/g, '').trim()

      const stars = ch('div.f6 a')
        .first()
        .text()
        .replace(/[\n\s,]+/g, '')

      const get_stars = ch('div.f6 span.float-sm-right')
        .last()
        .text()
        .replace(/[\n\s,a-zA-Z]+/g, '')

      const language = ch("span[itemprop='programmingLanguage']").text()

      trendRepos.push({
        owner_name,
        repo_name,
        html_url: `${this._baseUrl}/` + html_url,
        description,
        language,
        stargazers_count: parseInt(stars),
        get_stars: parseInt(get_stars),
      })
    }

    return trendRepos
  }

  async getTrendDevelopers(): Promise<ITrendDeveloper[]> {
    const res = await this._axios.get('/trending/developers')

    const html = await res.data

    const loader = load(html)

    const nodes = loader('article.Box-row')

    const trendDevelopers: ITrendDeveloper[] = []

    for (let i = 0; i < nodes.length; i++) {
      const ch = load(nodes[i])
      const avatar_url = ch('.mx-3 a img').attr('src') ?? ''
      const name = ch('div.col-md-6 h1.h3 a.Link').text().trim()
      const sub_name = ch('.col-md-6 p a').text().trim() ?? ''

      const repo_name = ch('.my-md-0 h1.lh-condensed a').text().trim()
      const html_url = ch('.my-md-0 h1.lh-condensed a').attr('href') ?? ''

      const description = ch('.my-md-0 .mt-1').text().replace(/\n+/g, '').trim()

      const url = html_url === '' ? sub_name : html_url

      trendDevelopers.push({
        avatar_url,
        name,
        sub_name,
        html_url: `${this._baseUrl}/` + url,
        repo_name,
        description,
      })
    }

    return trendDevelopers
  }
}

const client = new TrendsClient()

export default client
