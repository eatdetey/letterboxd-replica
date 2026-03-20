import { getMock } from '@shared/api/mockClient'
import type {
  MoviePreviewDto,
  UserListDto,
  UserProfileResponseDto,
  UserReviewDto,
} from '../model/types'

function safeArray<T>(value: T[] | undefined | null): T[] {
  return Array.isArray(value) ? value : []
}

async function tryGetMock<T>(path: string): Promise<T | null> {
  try {
    return await getMock<T>(path)
  } catch {
    return null
  }
}

function buildFallbackProfile(idOrUsername: string): UserProfileResponseDto {
  const username = idOrUsername || 'user'

  return {
    user: {
      id: idOrUsername || 'usr_mock',
      username,
      email: `${username}@example.com`,
      role: 'user',
      avatar_url: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=512&q=80',
      cover_url: 'https://images.unsplash.com/photo-1522202176988-66273c2fd55f?auto=format&fit=crop&w=1600&q=80',
      bio: 'Movie lover. Writes reviews and curates lists.',
      tags: ['Drama', 'Sci-Fi', 'Thriller'],
      stats: {
        filmsWatched: 0,
        reviewsCount: 0,
        listsCount: 0,
        averageRating: null,
        followersCount: 0,
        followingCount: 0,
        likesCount: 0,
      },
    },
    lists: [],
    reviews: [],
    recentRatings: [],
    watchlist: [],
  }
}

export const userProfileMockApi = {
  async getProfile(idOrUsername: string): Promise<UserProfileResponseDto> {
    const normalized = idOrUsername.trim()

    if (!normalized) {
      return buildFallbackProfile('usr_mock')
    }

    const base = `v1/users/${normalized}`

    const [profile, lists, reviews, recentRatings, watchlist] = await Promise.all([
      tryGetMock<UserProfileResponseDto>(`${base}/profile/response.json`),
      tryGetMock<{ items: UserListDto[] }>(`${base}/lists/response.json`),
      tryGetMock<{ items: UserReviewDto[] }>(`${base}/reviews/response.json`),
      tryGetMock<{ items: MoviePreviewDto[] }>(`${base}/recent-ratings/response.json`),
      tryGetMock<{ items: MoviePreviewDto[] }>(`${base}/watchlist/response.json`),
    ])

    if (profile) {
      return {
        user: profile.user,
        lists: safeArray(profile.lists),
        reviews: safeArray(profile.reviews),
        recentRatings: safeArray(profile.recentRatings),
        watchlist: safeArray(profile.watchlist),
      }
    }

    const fallback = buildFallbackProfile(normalized)

    return {
      user: fallback.user,
      lists: safeArray(lists?.items),
      reviews: safeArray(reviews?.items),
      recentRatings: safeArray(recentRatings?.items),
      watchlist: safeArray(watchlist?.items),
    }
  },
}