import requests
import uuid

def test_post_api_post_create_app_with_valid_data():
    base_url = "http://localhost:8080"
    timeout = 30

    # Step 1: Signup new user with dynamic email to get auth token
    unique_email = f"user_{uuid.uuid4()}@example.com"
    signup_payload = {
        "name": "Test User",
        "email": unique_email,
        "password": "ValidPass123!"
    }
    signup_response = requests.post(f"{base_url}/api/v1/auth/signup", json=signup_payload, timeout=timeout)
    assert signup_response.status_code == 201, f"Signup failed: {signup_response.text}"
    signup_data = signup_response.json()
    assert signup_data["status"] == "success"
    token = signup_data["data"]["token"]
    assert token, "Token not found in signup response"

    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }

    # Payload requirements: requestInterval as string to match server schema
    unique_app_name = f"My Test App {uuid.uuid4()}"
    app_payload = {
        "name": unique_app_name,
        "url": "https://my-test-app.example.com/health",
        "requestInterval": "10"
    }

    app_id = None
    try:
        # Step 2: POST /api/v1/apps/post to create a new app
        app_response = requests.post(f"{base_url}/api/v1/apps/post", json=app_payload, headers=headers, timeout=timeout)
        assert app_response.status_code == 201, f"App creation failed: {app_response.text}"
        app_response_data = app_response.json()
        assert app_response_data["status"] == "success"
        app_data = app_response_data["data"]["app"]

        # Validate the returned app data fields
        app_id = app_data.get("id")
        assert app_id, "App ID missing in response"
        assert app_data.get("name") == app_payload["name"]
        assert app_data.get("url") == app_payload["url"]
        assert app_data.get("requestInterval") == 10
        assert app_data.get("isEnabled") is False

    finally:
        # Cleanup: delete the created app if created
        if app_id:
            del_resp = requests.delete(f"{base_url}/api/v1/apps/delete/{app_id}", headers=headers, timeout=timeout)
            if del_resp.status_code not in (200, 204, 404):
                print(f"Warning: Unexpected status code deleting app {app_id}: {del_resp.status_code} {del_resp.text}")

test_post_api_post_create_app_with_valid_data()