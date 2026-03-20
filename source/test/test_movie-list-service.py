import pytest
import requests
import uuid
from typing import Dict, Any

BASE_URL = "https://localhost:8080/api/v1/playlists"
AUTH_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MywiVXNlcm5hbWUiOiJhc3lhbml4IiwiRW1haWwiOiJhc3lhQGdtYWlsLmNvbSIsIlN0YXR1cyI6ImFjdGl2ZSIsIlJvbGVzIjpbInVzZXIiXSwiSXNEZWxldGVkIjpmYWxzZSwic3ViIjoiMyIsImV4cCI6MTc3NDAzNDE3NywiaWF0IjoxNzc0MDMzMjc3fQ.XjrO5aT8bv6OxCwGeKAiEHCFDpGq-6RQnLLlvXkiBTM"

class TestPlaylistService:
    created_playlist_id = None
    playlist_name = f"My Favorites {uuid.uuid4().hex[:6]}"
    new_name = "Best Movies Ever"

    @pytest.fixture(autouse=True)
    def setup_headers(self):
        self.headers = {
            "Authorization": f"Bearer {AUTH_TOKEN}",
            "Content-Type": "application/json"
        }

    def test_1_create_playlist(self) -> None:
        """Создание нового плейлиста."""
        payload = {"name": self.playlist_name}
        response = requests.post(BASE_URL, json=payload, headers=self.headers)
        
        assert response.status_code in [200, 201], f"Failed to create playlist. Response: {response.text}"
        
        data = response.json()
        assert data.get("name") == self.playlist_name, "Playlist name mismatch."
        assert "id" in data, "ID is missing in response."
        
        TestPlaylistService.created_playlist_id = data["id"]

    def test_2_get_playlists_list(self) -> None:
        """Получение списка плейлистов и проверка наличия созданного."""
        response = requests.get(BASE_URL, headers=self.headers)
        
        assert response.status_code == 200
        data = response.json()
        
        assert "items" in data, "Field 'items' is missing."
        assert isinstance(data["items"], list), "'items' should be a list."
        
        playlist_ids = [item["id"] for item in data["items"]]
        assert self.created_playlist_id in playlist_ids, "Created playlist not found in the list."

    def test_3_rename_playlist(self) -> None:
        """Обновление плейлиста."""
        url = f"{BASE_URL}/{self.created_playlist_id}"
        payload = {"name": self.new_name}
        
        response = requests.put(url, json=payload, headers=self.headers)
        
        assert response.status_code == 200, f"Rename failed. Response: {response.text}"
        
        check_response = requests.get(BASE_URL, headers=self.headers)
        updated_item = next(item for item in check_response.json()["items"] 
                           if item["id"] == self.created_playlist_id)
        assert updated_item["name"] == self.new_name

    def test_4_delete_playlist(self) -> None:
        """Удаление плейлиста."""
        url = f"{BASE_URL}/{self.created_playlist_id}"
        
        response = requests.delete(url, headers=self.headers)
        
        assert response.status_code in [200, 204], f"Delete failed. Response: {response.text}"
        
        check_response = requests.get(BASE_URL, headers=self.headers)
        playlist_ids = [item["id"] for item in check_response.json()["items"]]
        assert self.created_playlist_id not in playlist_ids, "Playlist was not deleted."