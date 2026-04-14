import requests
import uuid

BASE_URL = "http://localhost:8080"
TIMEOUT = 30

def test_patch_apps_update_app_with_valid_changes():
    # Signup a new user to get auth token
    signup_email = f"user_{uuid.uuid4()}@example.com"
    signup_payload = {
        "name": "Test User",
        "email": signup_email,
        "password": "Secure@1234"
    }
    signup_resp = requests.post(f"{BASE_URL}/api/v1/auth/signup", json=signup_payload, timeout=TIMEOUT)
    assert signup_resp.status_code == 201, f"Signup failed: {signup_resp.text}"
    signup_data = signup_resp.json()
    assert signup_data.get("status") == "success", f"Unexpected signup status: {signup_data}"
    data = signup_data.get("data")
    assert data is not None, f"Missing 'data' in signup response: {signup_data}"
    user = data.get("user")
    assert user is not None, f"Missing 'user' in signup data: {data}"
    token = data.get("token")
    assert token is not None, f"Missing 'token' in signup data: {data}"
    user_id = user.get("id")
    assert user_id is not None, f"Missing 'id' in user data: {user}"
    headers = {"Authorization": f"Bearer {token}"}
    
    # Create a new app first (use requestInterval as integer on create per PRD)
    app_create_payload = {
        "name": "Original App Name",
        "url": "https://example.com/api",
        "requestInterval": 10
    }
    app_create_resp = requests.post(f"{BASE_URL}/api/v1/apps/post", json=app_create_payload, headers=headers, timeout=TIMEOUT)
    assert app_create_resp.status_code == 201, f"App creation failed: {app_create_resp.text}"
    app_create_data = app_create_resp.json()
    assert app_create_data.get("status") == "success", f"Unexpected app creation status: {app_create_data}"
    app = app_create_data["data"].get("app")
    assert app is not None, f"Missing 'app' in creation data: {app_create_data}" 
    app_id = app.get("id")
    assert app_id is not None, f"Missing 'id' in app data: {app}"

    try:
        # Prepare patch payload with updated name and requestInterval as integer 10
        patch_payload = {
            "name": "Updated App Name",
            "requestInterval": 10
        }
        patch_resp = requests.patch(f"{BASE_URL}/api/v1/apps/update/{app_id}", json=patch_payload, headers=headers, timeout=TIMEOUT)
        assert patch_resp.status_code == 200, f"App update failed: {patch_resp.text}"
        patch_data = patch_resp.json()
        assert patch_data.get("status") == "success", f"Unexpected patch status: {patch_data}"
        updated_app = patch_data["data"].get("app")
        assert updated_app is not None, f"Missing 'app' in patch response data: {patch_data}"
        assert updated_app.get("id") == app_id
        assert updated_app.get("name") == patch_payload["name"]
        ri = updated_app.get("requestInterval")
        # requestInterval should be integer 10
        assert isinstance(ri, int) and ri == 10, f"requestInterval is not correctly updated: {ri}"
    finally:
        # Clean up: delete the created app
        del_resp = requests.delete(f"{BASE_URL}/api/v1/apps/delete/{app_id}", headers=headers, timeout=TIMEOUT)
        # Either 200 success or 204 no content is acceptable; tolerate 404 for cleanup safety
        assert del_resp.status_code in [200, 204, 404]


test_patch_apps_update_app_with_valid_changes()
