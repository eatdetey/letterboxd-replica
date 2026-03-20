import pytest
import requests
from typing import Dict, Any

BASE_URL = "http://localhost:8080/api/v1/movies"
AUTH_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MywiVXNlcm5hbWUiOiJhc3lhbml4IiwiRW1haWwiOiJhc3lhQGdtYWlsLmNvbSIsIlN0YXR1cyI6ImFjdGl2ZSIsIlJvbGVzIjpbInVzZXIiXSwiSXNEZWxldGVkIjpmYWxzZSwic3ViIjoiMyIsImV4cCI6MTc3NDAzNDE3NywiaWF0IjoxNzc0MDMzMjc3fQ.XjrO5aT8bv6OxCwGeKAiEHCFDpGq-6RQnLLlvXkiBTM" 

class TestMovieService:
    
    @pytest.fixture
    def auth_headers(self) -> Dict[str, str]:
        return {"Authorization": f"Bearer {AUTH_TOKEN}"}

    def test_search_movies_by_title(self, auth_headers: Dict[str, str]) -> None:
        """
        Проверка поиска фильмов по названию.
        """
        params = {
            "search": "Военная машина",
            "limit": 10
        }
        response = requests.get(BASE_URL, params=params, headers=auth_headers)
        assert response.status_code == 200, f"Search failed. Response: {response.text}"
        
        data: Dict[str, Any] = response.json()
        assert "items" in data, "The 'items' field is missing in the response."
        assert isinstance(data["items"], list), f"Expected 'items' to be a list, but got {type(data['items'])}."
        
        assert "total" in data, "The 'total' field is missing in the response."
        assert isinstance(data["total"], int), f"Expected 'total' to be an integer, but got {type(data['total'])}."

    def test_get_movie_reviews(self, auth_headers: Dict[str, str]) -> None:
        """
        Проверка получения отзывов для конкретного фильма.
        Проверяем наличие списка items в ответе.
        """

        movie_id = "mov_001" 
        url = f"{BASE_URL}/{movie_id}/reviews"
        response = requests.get(url, headers=auth_headers)
        assert response.status_code == 200, f"Failed to get reviews. Response: {response.text}"
        
        data: Dict[str, Any] = response.json()
        assert "items" in data, "The 'items' field is missing in the reviews response."
        assert isinstance(data["items"], list), "Reviews 'items' should be a list."