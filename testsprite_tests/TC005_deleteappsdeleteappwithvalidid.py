import requests
import uuid

BASE_URL = "http://localhost:8080"
AUTH_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzYxODkwNTksImlhdCI6MTc3NjE1NjY1OSwidXNlcklkIjoiZGEwMGVkZDYtN2QyNS00MjA2LTlkZDAtZjY3MjZhOGU0ZWEzIn0.ofgxmJeMaAcnvIqrE5qYXzmD3ayJGQTpGTSatgwaOSs"
HEADERS = {"Authorization": f"Bearer {AUTH_TOKEN}"}
TIMEOUT = 30

def test_deleteappsdeleteappwithvalidid():
    # Step 1: Create a new app to delete later
    unique_email = f"user_{uuid.uuid4()}@example.com"
    signup_payload = {
        "name": "Delete Test User",
        "email": unique_email,
        "password": "Secure@1234"
    }
    # Signup new user
    signup_resp = requests.post(f"{BASE_URL}/api/v1/auth/signup", json=signup_payload, timeout=TIMEOUT)
    assert signup_resp.status_code == 201, f"Signup failed: {signup_resp.text}"
    signup_data = signup_resp.json()
    assert signup_data.get("status") == "success"
    assert isinstance(signup_data.get("data"), dict), "Missing data object in signup response"
    assert "user" in signup_data["data"], "User data missing in signup response"
    assert "token" in signup_data["data"], "Token missing in signup response"
    user_token = signup_data["data"]["token"]
    user_id = signup_data["data"]["user"]["id"]
    user_headers = {"Authorization": f"Bearer {user_token}"}

    # Create an app under this user
    app_payload = {
        "name": "App to Delete",
        "url": "https://example-app-to-delete.com",
        "requestInterval": 10
    }
    create_app_resp = requests.post(f"{BASE_URL}/api/v1/apps/post", json=app_payload, headers=user_headers, timeout=TIMEOUT)
    assert create_app_resp.status_code == 201, f"App creation failed: {create_app_resp.text}"
    app_data = create_app_resp.json()
    assert app_data.get("status") == "success"
    assert "app" in app_data.get("data", {}), "App data missing in create app response"
    app_id = app_data["data"]["app"]["id"]

    try:
        # Step 2: DELETE the created app with valid ID
        delete_resp = requests.delete(f"{BASE_URL}/api/v1/apps/delete/{app_id}", headers=user_headers, timeout=TIMEOUT)
        assert delete_resp.status_code in (200, 204), f"Valid delete failed: {delete_resp.text}"
        if delete_resp.status_code != 204:
            delete_json = delete_resp.json()
            assert delete_json.get("status") == "success"

        # Step 3: Attempt to DELETE a non-existent app ID and expect 404
        non_existent_app_id = "00000000-0000-0000-0000-000000000000"
        delete_nonexist_resp = requests.delete(f"{BASE_URL}/api/v1/apps/delete/{non_existent_app_id}", headers=user_headers, timeout=TIMEOUT)
        assert delete_nonexist_resp.status_code == 404, f"Delete of non-existent app did not return 404: {delete_nonexist_resp.text}"
        err_json = delete_nonexist_resp.json()
        assert err_json.get("status") == "error"
        assert "app not found" in err_json.get("message", "").lower()
    finally:
        # Cleanup: Ensure the created app is deleted if it still exists
        # Check if app still exists by trying to delete again quietly
        try:
            requests.delete(f"{BASE_URL}/api/v1/apps/delete/{app_id}", headers=user_headers, timeout=TIMEOUT)
        except:
            pass

test_deleteappsdeleteappwithvalidid()
