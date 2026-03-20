import uuid
import pytest
import requests
from typing import Dict, Any

BASE_URL = "http://localhost:8080/api/v1/auth"

@pytest.fixture(scope="class")
def user_data() -> Dict[str, str]:
    unique_suffix = str(uuid.uuid4())[:8]
    return {
        "username": f"user_{unique_suffix}",
        "password": "StrongPassword123!",
        "email": f"user_{unique_suffix}@example.com"
    }

class TestUserAuthJourney:
    def test_successful_registration(self, user_data: Dict[str, str]) -> None:
        """Тест успешной регистрации нового пользователя."""
        
        response = requests.post(f"{BASE_URL}/register", json=user_data)
        assert response.status_code == 200, f"Expected 200 Created, got {response.status_code}. Response: {response.text}"
        data: Dict[str, Any] = response.json()

        assert data.get("user").get("username") == user_data["username"], "Username in response does not match the requested one."
        assert data.get("user").get("email") == user_data["email"], "Email in response does not match the requested one."
        assert data.get("user").get("role") == "ROLE_USER", "Default role should be 'ROLE_USER'."
        assert data.get("user").get("status") == "active", "Default status should be 'active'."
        
        assert "id" in data.get("user"), "User ID is missing in response."
        assert "access_token" in data, "Access token is missing in response."

    def test_unsuccessful_registration_duplicate(self, user_data: Dict[str, str]) -> None:
        """Тест неуспешной регистрации."""
        
        response = requests.post(f"{BASE_URL}/register", json=user_data)
        assert response.status_code in [400, 409], f"Expected 400 or 409 for duplicate registration, got {response.status_code}."

    def test_successful_login(self, user_data: Dict[str, str]) -> None:
        """Тест успешного входа, используя креды пользователя, созданного в первом тесте."""
        
        login_payload = {
            "username": user_data["username"],
            "password": user_data["password"]
        }
        response = requests.post(f"{BASE_URL}/login", json=login_payload)
        assert response.status_code == 200, f"Expected 200 OK, got {response.status_code}. Response: {response.text}"
        
        data: Dict[str, Any] = response.json()
        assert data.get("user").get("username") == user_data["username"], "Logged in username mismatch."
        assert "id" in data.get("user"), "User ID is missing in login response."
        assert "access_token" in data, "Access token is missing in response."
