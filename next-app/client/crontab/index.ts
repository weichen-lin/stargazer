'use server'

import BaseClient from '../base-client'

export interface ICrontab {
  triggered_at: string
  created_at: string
  updated_at: string
  status: string
  last_triggered_at: string
  version: number
}

class CrontabClient extends BaseClient {
  getCrontab() {
    return this.get<ICrontab>('/crontab')
  }

  createCronTab() {
    return this.post<any, ICrontab>('/crontab', null)
  }

  updateCronTab(triggered_at: string) {
    return this.patch(`/crontab?triggered_at=${triggered_at}`, null)
  }
}

export default CrontabClient
