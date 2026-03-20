const ACCESS_TOKEN_KEY = 'cinevault_access_token'
const REFRESH_TOKEN_KEY = 'cinevault_refresh_token'
const USER_KEY = 'cinevault_user'

export const authTokenStorage = {
  getAccessToken(): string | null {
    return localStorage.getItem(ACCESS_TOKEN_KEY)
  },

  getRefreshToken(): string | null {
    return localStorage.getItem(REFRESH_TOKEN_KEY)
  },

  setSession(user: unknown, accessToken: string, refreshToken: string) {
    localStorage.setItem(USER_KEY, JSON.stringify(user))
    localStorage.setItem(ACCESS_TOKEN_KEY, accessToken)
    localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)
  },

  clearSession() {
    localStorage.removeItem(USER_KEY)
    localStorage.removeItem(ACCESS_TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
  },

  getUser<T>(): T | null {
    const raw = localStorage.getItem(USER_KEY)
    if (!raw) return null

    try {
      return JSON.parse(raw) as T
    } catch {
      return null
    }
  },
}