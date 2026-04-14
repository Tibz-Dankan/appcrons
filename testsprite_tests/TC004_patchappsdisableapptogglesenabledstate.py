import requests
import uuid

BASE_URL = "http://localhost:8080"
TIMEOUT = 30

def test_patch_apps_disable_app_toggles_enabled_state():
    # Step 1: Signup a new user and get accessToken
    signup_url = f"{BASE_URL}/api/v1/auth/signup"
    unique_email = f"{uuid.uuid4()}@test.com"
    signup_payload = {
        "name": "Test User",
        "email": unique_email,
        "password": "TestPass@1234"
    }
    signup_resp = requests.post(signup_url, json=signup_payload, timeout=TIMEOUT)
    assert signup_resp.status_code == 201, f"Signup failed: {signup_resp.text}"
    signup_data = signup_resp.json()
    access_token = signup_data.get("accessToken")
    assert access_token, "No accessToken returned from signup"
    headers = {"Authorization": f"Bearer {access_token}"}

    # Step 2: Create a new app; new apps are disabled by default (isDisabled = True)
    create_app_url = f"{BASE_URL}/api/v1/apps/post"
    app_name = f"App {uuid.uuid4()}"
    app_url = f"https://app-{uuid.uuid4()}.onrender.com/active"
    create_app_payload = {
        "name": app_name,
        "url": app_url,
        "requestInterval": "10"
    }
    create_app_resp = requests.post(create_app_url, json=create_app_payload, headers=headers, timeout=TIMEOUT)
    assert create_app_resp.status_code == 201, f"App creation failed: {create_app_resp.text}"
    app_data = create_app_resp.json()
    assert app_data.get("status") == "success"
    app = app_data.get("data", {}).get("app")
    assert app, "No app object in response"
    app_id = app.get("id")
    assert app_id, "No app ID returned"

    try:
        # Step 3: Enable the app first (PATCH /api/v1/apps/enable/{appId})
        enable_url = f"{BASE_URL}/api/v1/apps/enable/{app_id}"
        enable_resp = requests.patch(enable_url, headers=headers, timeout=TIMEOUT)
        assert enable_resp.status_code == 200, f"Enabling app failed: {enable_resp.text}"
        enable_resp_data = enable_resp.json()
        assert enable_resp_data.get("status") == "success"
        enabled_app = enable_resp_data.get("data", {}).get("app")
        assert enabled_app is not None
        assert enabled_app.get("isDisabled") is False, "App should be enabled (isDisabled=False)"

        # Step 4: Disable the enabled app (PATCH /api/v1/apps/disable/{appId})
        disable_url = f"{BASE_URL}/api/v1/apps/disable/{app_id}"
        disable_resp = requests.patch(disable_url, headers=headers, timeout=TIMEOUT)
        assert disable_resp.status_code == 200, f"Disabling app failed: {disable_resp.text}"
        disable_resp_data = disable_resp.json()
        assert disable_resp_data.get("status") == "success"
        disabled_app = disable_resp_data.get("data", {}).get("app")
        assert disabled_app is not None
        assert disabled_app.get("isDisabled") is True, "App should be disabled (isDisabled=True)"

        # Step 5: Attempt to disable again (should fail with 400 and message "app is already disabled")
        disable_again_resp = requests.patch(disable_url, headers=headers, timeout=TIMEOUT)
        assert disable_again_resp.status_code == 400, f"Expected 400 when disabling already disabled app, got {disable_again_resp.status_code}"
        disable_again_data = disable_again_resp.json()
        assert disable_again_data.get("status") == "error"
        assert disable_again_data.get("message") == "app is already disabled"

    finally:
        # Cleanup: Delete the app
        delete_url = f"{BASE_URL}/api/v1/apps/delete/{app_id}"
        delete_resp = requests.delete(delete_url, headers=headers, timeout=TIMEOUT)
        assert delete_resp.status_code == 200, f"App deletion failed: {delete_resp.text}"

test_patch_apps_disable_app_toggles_enabled_state()