import { env } from '@shared/config/env'
import { authTokenStorage } from './authTokenStorage'

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE'
type RequestOptions = RequestInit & { skipAuthRefresh?: boolean }

export class HttpClient {
  private readonly baseUrl: string
  private refreshPromise: Promise<string | null> | null = null

  constructor(baseUrl = env.apiBaseUrl) {
    this.baseUrl = baseUrl
  }

  private isAuthEndpoint(path: string): boolean {
    return path.startsWith('/v1/auth/')
  }

  private hasLocalSession(): boolean {
    return Boolean(authTokenStorage.getAccessToken()) || Boolean(authTokenStorage.getUser<unknown>())
  }

  private async refreshAccessToken(): Promise<string | null> {
    if (this.refreshPromise) {
      return this.refreshPromise
    }

    this.refreshPromise = (async () => {
      try {
        const response = await fetch(`${this.baseUrl}/v1/auth/refresh`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
        })

        if (!response.ok) {
          return null
        }

        const payload = await response.json() as { access_token?: string }
        const accessToken = payload.access_token
        if (!accessToken) {
          return null
        }

        const user = authTokenStorage.getUser<unknown>()
        if (user) {
          authTokenStorage.setSession(user, accessToken)
        } else {
          authTokenStorage.setAccessToken(accessToken)
        }

        return accessToken
      } catch {
        return null
      } finally {
        this.refreshPromise = null
      }
    })()

    return this.refreshPromise
  }

  async request<T>(path: string, options: RequestOptions = {}): Promise<T> {
    const { skipAuthRefresh = false, ...fetchOptions } = options
    const headers = new Headers(options.headers)

    if (!headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json')
    }

    const accessToken = authTokenStorage.getAccessToken()
    if (accessToken) {
      headers.set('Authorization', `Bearer ${accessToken}`)
    }

    const response = await fetch(`${this.baseUrl}${path}`, {
      ...fetchOptions,
      headers,
      credentials: 'include',
    })

    if (!response.ok) {
      if (response.status === 401 && !skipAuthRefresh && !this.isAuthEndpoint(path)) {
        const refreshedToken = await this.refreshAccessToken()
        if (refreshedToken) {
          return this.request<T>(path, {
            ...fetchOptions,
            skipAuthRefresh: true,
          })
        }
      }

      if (response.status === 401 && !this.isAuthEndpoint(path) && this.hasLocalSession()) {
        authTokenStorage.clearSession()
        window.location.assign('/login')
      }

      const errorText = await response.text()
      throw new Error(`HTTP ${response.status}: ${errorText}`)
    }

    if (response.status === 204) {
      return undefined as T
    }

    const contentLength = response.headers.get('content-length')
    if (contentLength === '0') {
      return undefined as T
    }

    return response.json() as Promise<T>
  }

  get<T>(path: string): Promise<T> {
    return this.request<T>(path)
  }

  post<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>(path, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  put<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>(path, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  delete<T>(path: string): Promise<T> {
    return this.request<T>(path, { method: 'DELETE' })
  }
}

export const httpClient = new HttpClient()
