import requests
import uuid

BASE_URL = "http://localhost:8080"
TIMEOUT = 30

def test_patchappsupdateappwithvalidchanges():
    # Step 1: Signup new user to get fresh accessToken
    signup_url = f"{BASE_URL}/api/v1/auth/signup"
    unique_email = f"{uuid.uuid4()}@test.com"
    signup_payload = {
        "name": "Test User",
        "email": unique_email,
        "password": "TestPass@1234"
    }
    signup_resp = requests.post(signup_url, json=signup_payload, timeout=TIMEOUT)
    assert signup_resp.status_code == 201, f"Signup failed: {signup_resp.text}"
    signup_json = signup_resp.json()
    assert "accessToken" in signup_json, "No accessToken in signup response"
    token = signup_json["accessToken"]
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }

    # Step 2: Create an app to update later
    create_app_url = f"{BASE_URL}/api/v1/apps/post"
    app_name_orig = f"App {uuid.uuid4()}"
    app_url_orig = f"https://app-{uuid.uuid4()}.onrender.com/active"
    create_app_payload = {
        "name": app_name_orig,
        "url": app_url_orig,
        "requestInterval": "10"
    }
    create_resp = requests.post(create_app_url, headers=headers, json=create_app_payload, timeout=TIMEOUT)
    assert create_resp.status_code == 201, f"App creation failed: {create_resp.text}"
    create_json = create_resp.json()
    assert create_json.get("status") == "success", f"Unexpected status on app creation: {create_json}"
    app = create_json.get("data", {}).get("app")
    assert app and "id" in app, "App ID missing in creation response"
    app_id = app["id"]

    try:
        # Step 3: PATCH to update app details: change name and requestInterval
        update_app_url = f"{BASE_URL}/api/v1/apps/update/{app_id}"
        updated_name = f"Updated {uuid.uuid4()}"
        updated_url = f"https://upd-{uuid.uuid4()}.onrender.com/active"
        update_payload = {
            "name": updated_name,
            "url": updated_url,
            "requestInterval": "15"
        }
        patch_resp = requests.patch(update_app_url, headers=headers, json=update_payload, timeout=TIMEOUT)
        assert patch_resp.status_code == 200, f"App update failed: {patch_resp.text}"
        patch_json = patch_resp.json()
        assert patch_json.get("status") == "success", f"Unexpected status on app update: {patch_json}"
        updated_app = patch_json.get("data", {}).get("app")
        assert updated_app, "Updated app data missing in response"

        # Validate updated fields
        assert updated_app.get("id") == app_id, "Updated app ID mismatch"
        assert updated_app.get("name") == updated_name, "App name not updated correctly"
        assert updated_app.get("url") == updated_url, "App url not updated correctly"
        assert updated_app.get("requestInterval") == "15", "requestInterval not updated correctly"

    finally:
        # Cleanup: delete the created app
        delete_app_url = f"{BASE_URL}/api/v1/apps/delete/{app_id}"
        del_resp = requests.delete(delete_app_url, headers=headers, timeout=TIMEOUT)
        assert del_resp.status_code == 200, f"App deletion failed during cleanup: {del_resp.text}"

test_patchappsupdateappwithvalidchanges()