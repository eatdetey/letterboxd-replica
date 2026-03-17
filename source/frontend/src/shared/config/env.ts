const useMocksEnv = import.meta.env.VITE_USE_MOCKS

export const env = {
  useMocks: useMocksEnv ? useMocksEnv === 'true' : true,
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || '/api',
}
