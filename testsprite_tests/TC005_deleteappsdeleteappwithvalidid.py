import requests
import uuid

BASE_URL = "http://localhost:8080"
TIMEOUT = 30


def test_deleteappsdeleteappwithvalidid():
    # Step 1: Signup to get fresh user and accessToken
    signup_url = f"{BASE_URL}/api/v1/auth/signup"
    email = f"{uuid.uuid4()}@test.com"
    signup_payload = {
        "name": "Test User",
        "email": email,
        "password": "TestPass@1234"
    }
    signup_resp = requests.post(signup_url, json=signup_payload, timeout=TIMEOUT)
    assert signup_resp.status_code == 201, f"Signup failed: {signup_resp.text}"
    signup_json = signup_resp.json()
    assert "accessToken" in signup_json, "No accessToken in signup response"
    access_token = signup_json["accessToken"]
    headers = {
        "Authorization": f"Bearer {access_token}"
    }

    # Step 2: Create an app as fixture
    create_app_url = f"{BASE_URL}/api/v1/apps/post"
    app_name = f"App {uuid.uuid4()}"
    app_url = f"https://app-{uuid.uuid4()}.onrender.com/active"
    create_app_payload = {
        "name": app_name,
        "url": app_url,
        "requestInterval": "10"
    }
    create_resp = requests.post(create_app_url, json=create_app_payload, headers=headers, timeout=TIMEOUT)
    assert create_resp.status_code == 201, f"App creation failed: {create_resp.text}"
    create_json = create_resp.json()
    assert create_json.get("status") == "success", f"App creation status not success: {create_resp.text}"
    app_data = create_json.get("data", {}).get("app")
    assert app_data and "id" in app_data, "No app id in response"
    app_id = app_data["id"]

    try:
        # Step 3: Delete the created app - expect 200 success
        delete_url = f"{BASE_URL}/api/v1/apps/delete/{app_id}"
        delete_resp = requests.delete(delete_url, headers=headers, timeout=TIMEOUT)
        assert delete_resp.status_code == 200, f"Failed to delete existing app: {delete_resp.text}"

        # Step 4: Attempt to delete the same app again - expect 404 Not Found since app was deleted
        delete_resp_404 = requests.delete(delete_url, headers=headers, timeout=TIMEOUT)
        assert delete_resp_404.status_code == 404, f"Deleting non-existent app did not return 404: {delete_resp_404.text}"
    finally:
        # Cleanup: Just in case deletion failed, attempt to delete app to avoid leftover state
        requests.delete(f"{BASE_URL}/api/v1/apps/delete/{app_id}", headers=headers, timeout=TIMEOUT)


test_deleteappsdeleteappwithvalidid()