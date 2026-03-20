const useMocksEnv = import.meta.env.VITE_USE_MOCKS
const useMoviesMocksEnv = import.meta.env.VITE_USE_MOCKS_MOVIES
const useReviewsMocksEnv = import.meta.env.VITE_USE_MOCKS_REVIEWS

function resolveUseMocks(value: string | undefined, fallback: boolean) {
  return value ? value === 'true' : fallback
}

const useMocks = resolveUseMocks(useMocksEnv, true)

export const env = {
  useMocks,
  useMovieMocks: resolveUseMocks(useMoviesMocksEnv, useMocks),
  useReviewMocks: resolveUseMocks(useReviewsMocksEnv, useMocks),
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || '/api',
}
