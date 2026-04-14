import requests
import uuid

BASE_URL = "http://localhost:8080"
TIMEOUT = 30
AUTH_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzYxODkwNTksImlhdCI6MTc3NjE1NjY1OSwidXNlcklkIjoiZGEwMGVkZDYtN2QyNS00MjA2LTlkZDAtZjY3MjZhOGU0ZWEzIn0.ofgxmJeMaAcnvIqrE5qYXzmD3ayJGQTpGTSatgwaOSs"
HEADERS_AUTH = {"Authorization": f"Bearer {AUTH_TOKEN}"}
HEADERS_JSON = {"Content-Type": "application/json"}

def test_patch_apps_disable_app_toggles_enabled_state():
    # 1. Signup user to get unique email and token
    signup_email = f"user_{uuid.uuid4()}@example.com"
    signup_payload = {
        "name": "Test User",
        "email": signup_email,
        "password": "SecurePass123!"
    }
    signup_resp = requests.post(f"{BASE_URL}/api/v1/auth/signup", json=signup_payload, timeout=TIMEOUT)
    assert signup_resp.status_code == 201, f"Signup failed with status {signup_resp.status_code}: {signup_resp.text}"
    signup_data = signup_resp.json()
    assert signup_data["status"] == "success"
    user_token = signup_data["data"]["token"]
    user_headers = {"Authorization": f"Bearer {user_token}", "Content-Type": "application/json"}

    # 2. Create a new app (isEnabled should be false initially)
    unique_url = f"https://{uuid.uuid4()}.example.com/health"
    app_payload = {
        "name": "App to Disable Test",
        "url": unique_url,
        "requestInterval": "10"  # must be string as per API error and PRD
    }
    app_create_resp = requests.post(f"{BASE_URL}/api/v1/apps/post", headers=user_headers, json=app_payload, timeout=TIMEOUT)
    assert app_create_resp.status_code == 201, f"App creation failed: {app_create_resp.text}"
    app_create_data = app_create_resp.json()
    assert app_create_data["status"] == "success"
    app = app_create_data["data"]["app"]
    app_id = app["id"]

    try:
        # 3. Enable the app first to ensure it is enabled before disabling
        enable_resp = requests.patch(f"{BASE_URL}/api/v1/apps/enable/{app_id}", headers=user_headers, timeout=TIMEOUT)
        assert enable_resp.status_code == 200, f"Enable app failed: {enable_resp.text}"
        enable_data = enable_resp.json()
        assert enable_data["status"] == "success"
        assert enable_data["data"]["isEnabled"] is True

        # 4. PATCH disable the enabled app
        disable_resp = requests.patch(f"{BASE_URL}/api/v1/apps/disable/{app_id}", headers=user_headers, timeout=TIMEOUT)
        assert disable_resp.status_code == 200, f"Disable app failed: {disable_resp.text}"
        disable_data = disable_resp.json()
        assert disable_data["status"] == "success"
        assert disable_data["data"]["isEnabled"] is False

        # 5. PATCH disable the already disabled app (idempotency)
        disable_again_resp = requests.patch(f"{BASE_URL}/api/v1/apps/disable/{app_id}", headers=user_headers, timeout=TIMEOUT)
        assert disable_again_resp.status_code == 200, f"Second disable failed: {disable_again_resp.text}"
        disable_again_data = disable_again_resp.json()
        assert disable_again_data["status"] == "success"
        assert disable_again_data["data"]["isEnabled"] is False

    finally:
        # Cleanup: delete the app
        del_resp = requests.delete(f"{BASE_URL}/api/v1/apps/delete/{app_id}", headers=user_headers, timeout=TIMEOUT)
        # Accept 200/204 success or 404 if already deleted
        assert del_resp.status_code in (200, 204, 404), f"App deletion failed: {del_resp.status_code}, {del_resp.text}"

test_patch_apps_disable_app_toggles_enabled_state()