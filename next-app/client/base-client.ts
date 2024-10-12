'use server'

import { sign } from 'jsonwebtoken'
import { z } from 'zod'
import { randomUUID } from 'crypto'
import axios, { AxiosInstance, AxiosRequestConfig, isAxiosError } from 'axios'
import { getServerSession } from 'next-auth'
import { options } from '@/app/api/auth/[...nextauth]/option'

const userSchema = z.object({
  name: z.string().min(1),
  email: z.string().email(),
  image: z.string(),
})

interface APIResponse<T> {
  status_code: number
  data: T
}

class BaseClient {
  private _health: boolean = false
  private readonly _secret: string
  private readonly _baseApiPath: string
  private readonly _options = options
  private headers: Record<string, string> = {}
  private _axiosInstance!: AxiosInstance

  constructor() {
    this._baseApiPath = process.env.STARGAZER_BACKEND as string
    this._secret = process.env.JWT_SECRET as string

    this._health = this.validate()

    this.appendDefaultHeaders({
      'X-Request-Id': randomUUID(),
    })

    this._axiosInstance = axios.create({
      baseURL: this._baseApiPath,
      headers: {
        'Content-Type': 'application/json',
        ...this.headers,
      },
    })
  }

  health(): boolean {
    return this._health
  }

  private validate(): boolean {
    const { success: baseUrlValidate } = z.string().url().safeParse(this._baseApiPath)
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

  private async api<T>(options: AxiosRequestConfig): Promise<APIResponse<T>> {
    const { url, method, headers, params, data, signal } = options

    if (!this._health) {
      return {
        status_code: 500,
        data: 'Internal Server Error' as T,
      }
    }

    const session = await getServerSession(this._options)
    if (!session) {
      return {
        status_code: 401,
        data: 'Unauthorized' as T,
      }
    }

    const { email } = userSchema.parse(session?.user)

    this.appendDefaultHeaders({
      Authorization: this.generateAccessToken(email),
      'X-User-Email': email,
    })

    try {
      const response = await this._axiosInstance.request<T>({
        url,
        method,
        headers: {
          ...this.headers,
          ...headers,
        },
        params,
        data,
        signal,
      })

      return {
        status_code: response.status,
        data: response.data,
      }
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response) {
          return {
            status_code: error.response.status,
            data: error.response.data,
          }
        }
      }

      return {
        status_code: 500,
        data: error as T,
      }
    }
  }

  protected async get<T>(url: string, params?: URLSearchParams, signal?: AbortSignal) {
    return this.api<T>({ url, method: 'GET', params, signal })
  }

  protected async post<V, T>(url: string, data: V, signal?: AbortSignal) {
    return this.api<T>({ url, method: 'POST', data, signal })
  }

  protected async patch<V, T>(url: string, data: V, signal?: AbortSignal) {
    return this.api<T>({ url, method: 'PATCH', data, signal })
  }

  protected async delete<V, T>(url: string, data?: V, signal?: AbortSignal) {
    return this.api<T>({ url, method: 'DELETE', data, signal })
  }

  protected async put<V, T>(url: string, data: V, signal?: AbortSignal) {
    return this.api<T>({ url, method: 'PUT', data, signal })
  }

  protected async head<T>(url: string, signal?: AbortSignal) {
    return this.api<T>({ url, method: 'HEAD', signal })
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
