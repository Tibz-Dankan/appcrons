import requests
import uuid

BASE_URL = "http://localhost:8080"
TIMEOUT = 30


def test_patch_apps_enable_app_toggles_enabled_state():
    # Step 1: Signup a new user to get fresh accessToken
    signup_url = f"{BASE_URL}/api/v1/auth/signup"
    unique_email = f"{uuid.uuid4()}@test.com"
    signup_payload = {
        "name": "Test User",
        "email": unique_email,
        "password": "TestPass@1234"
    }
    signup_resp = requests.post(signup_url, json=signup_payload, timeout=TIMEOUT)
    assert signup_resp.status_code == 201, f"Signup failed with status {signup_resp.status_code}"
    signup_json = signup_resp.json()
    assert "accessToken" in signup_json, "No accessToken in signup response"
    access_token = signup_json["accessToken"]
    headers = {"Authorization": f"Bearer {access_token}"}

    # Step 2: Create a new app (will be disabled by default)
    create_app_url = f"{BASE_URL}/api/v1/apps/post"
    app_name = f"App {uuid.uuid4()}"
    app_url_val = f"https://app-{uuid.uuid4()}.onrender.com/active"
    create_app_payload = {
        "name": app_name,
        "url": app_url_val,
        "requestInterval": "10"
    }
    create_resp = requests.post(create_app_url, json=create_app_payload, headers=headers, timeout=TIMEOUT)
    assert create_resp.status_code == 201, f"App creation failed with status {create_resp.status_code}"
    create_json = create_resp.json()
    assert create_json.get("status") == "success", "App creation response status not success"
    app = create_json.get("data", {}).get("app")
    assert app is not None, "No app data in response"
    app_id = app.get("id")
    assert app.get("isDisabled") is True, "New app should be disabled by default"

    enable_url = f"{BASE_URL}/api/v1/apps/enable/{app_id}"

    try:
        # Step 3: Enable the disabled app
        enable_resp_1 = requests.patch(enable_url, headers=headers, timeout=TIMEOUT)
        assert enable_resp_1.status_code == 200, f"Enable app failed with status {enable_resp_1.status_code}"
        enable_json_1 = enable_resp_1.json()
        assert enable_json_1.get("status") == "success", "Enable app response status not success"
        app_enabled_1 = enable_json_1.get("data", {}).get("app")
        assert app_enabled_1 is not None, "No app data in enable response"
        assert app_enabled_1.get("isDisabled") is False, "App should be enabled (isDisabled=False)"

        # Step 4: Enable the already enabled app again to test idempotency
        enable_resp_2 = requests.patch(enable_url, headers=headers, timeout=TIMEOUT)
        assert enable_resp_2.status_code == 200, f"Enable already enabled app failed with status {enable_resp_2.status_code}"
        enable_json_2 = enable_resp_2.json()
        assert enable_json_2.get("status") == "success", "Enable already enabled app response status not success"
        app_enabled_2 = enable_json_2.get("data", {}).get("app")
        assert app_enabled_2 is not None, "No app data in second enable response"
        assert app_enabled_2.get("isDisabled") is False, "App should remain enabled (isDisabled=False) after second enable call"

    finally:
        # Cleanup: Delete the created app
        delete_url = f"{BASE_URL}/api/v1/apps/delete/{app_id}"
        delete_resp = requests.delete(delete_url, headers=headers, timeout=TIMEOUT)
        assert delete_resp.status_code == 200, f"App deletion failed with status {delete_resp.status_code}"


test_patch_apps_enable_app_toggles_enabled_state()