'use server'

import BaseClient from './base-client'

export interface ICrontab {
  triggered_at: string
  created_at: string
  updated_at: string
  status: string
  last_triggered_at: string
  version: number
}

class CrontabClient extends BaseClient {
  constructor(email: string) {
    super(email)
  }

  getCrontab() {
    return this.get<ICrontab>('/crontab')
  }

  createCronTab() {
    return this.post<any, ICrontab>('/crontab', null)
  }

  updateCronTab(hour: string) {
    return this.patch(`/crontab?hour=${hour}`, null)
  }
}

export default CrontabClient
