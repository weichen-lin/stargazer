'use server'

import { sign } from 'jsonwebtoken'
import { z } from 'zod'
import { randomUUID } from 'crypto'

class BaseClient {
  private _health: boolean = false
  private readonly _secret!: string
  private readonly baseApiPath: string
  private headers: Record<string, string> = {}
  constructor(email: string) {
    this.baseApiPath = process.env.STARGAZER_BACKEND as string
    this._secret = process.env.JWT_SECRET as string

    this._health = this.validate()

    this.appendDefaultHeaders({
      Authorization: this.generateAccessToken(email),
      'X-Request-Id': randomUUID(),
      'X-User-Email': email,
    })
  }

  health(): boolean {
    return this._health
  }

  private validate(): boolean {
    const { success: baseUrlValidate } = z.string().url().safeParse(this.baseApiPath)
    const { success: jwtSecretValidate } = z.string().min(16).safeParse(this._secret)

    return baseUrlValidate && jwtSecretValidate
  }

  private generateAccessToken(email: string) {
    const expired = new Date()
    expired.setMinutes(expired.getMinutes() + 1)

    const payload = {
      email: email,
      create_at: new Date(),
      expire_at: expired,
    }

    const token = sign(payload, this._secret, { algorithm: 'HS256' })

    return token
  }

  private async api<T>(path: string, options: RequestInit): Promise<T> {
    const url = this.baseApiPath + path
    const { method, headers: extraHeaders, body, ...leftOptions } = options

    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        ...this.headers,
        ...extraHeaders,
      },
      body,
      ...leftOptions,
    })
    const data = await response.json()
    return data
  }

  protected async get<T>(url: string, signal?: AbortSignal) {
    return this.api<T>(url, { method: 'GET', signal })
  }

  protected async post<T>(url: string, body: any, signal?: AbortSignal) {
    return this.api<T>(url, { method: 'POST', body, signal })
  }

  protected async patch<T>(url: string, body: any, signal?: AbortSignal) {
    return this.api<T>(url, { method: 'PATCH', body, signal })
  }

  protected async delete<T>(url: string, signal?: AbortSignal) {
    return this.api<T>(url, { method: 'DELETE', signal })
  }

  protected async put<T>(url: string, body: any, signal?: AbortSignal) {
    return this.api<T>(url, { method: 'PUT', body, signal })
  }

  protected async head<T>(url: string, signal?: AbortSignal) {
    return this.api<T>(url, { method: 'HEAD', signal })
  }

  protected async options<T>(url: string, signal?: AbortSignal) {
    return this.api<T>(url, { method: 'OPTIONS', signal })
  }

  protected async connect<T>(url: string, signal?: AbortSignal) {
    return this.api<T>(url, { method: 'CONNECT', signal })
  }

  protected async trace<T>(url: string, signal?: AbortSignal) {
    return this.api<T>(url, { method: 'TRACE', signal })
  }

  public appendDefaultHeaders(headers: Record<string, string>) {
    if (!headers) return
    this.headers = {
      ...this.headers,
      ...headers,
    }
  }
}

export default BaseClient
