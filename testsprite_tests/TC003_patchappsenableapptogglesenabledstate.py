import requests
import uuid

BASE_URL = "http://localhost:8080"
TIMEOUT = 30

def test_patchappsenableapptogglesenabledstate():
    # Step 1: Signup new user with unique email
    signup_url = f"{BASE_URL}/api/v1/auth/signup"
    unique_email = f"user_{uuid.uuid4()}@example.com"
    signup_payload = {
        "name": "Test User",
        "email": unique_email,
        "password": "Secure@1234"
    }
    signup_resp = requests.post(signup_url, json=signup_payload, timeout=TIMEOUT)
    assert signup_resp.status_code == 201, f"Signup failed: {signup_resp.text}"
    signup_data = signup_resp.json()
    assert signup_data.get("status") == "success"
    token = signup_data["data"]["token"]
    user_id = signup_data["data"]["user"]["id"]

    headers = {
        "Authorization": f"Bearer {token}"
    }

    # Step 2: Create a new app which will be disabled by default
    create_app_url = f"{BASE_URL}/api/v1/apps/post"
    app_payload = {
        "name": "Test App for Enable Endpoint",
        "url": "https://example.com",
        "requestInterval": "10"  # Must be string as per instructions
    }
    create_resp = requests.post(create_app_url, json=app_payload, headers=headers, timeout=TIMEOUT)
    assert create_resp.status_code == 201, f"App creation failed: {create_resp.text}"
    create_data = create_resp.json()
    assert create_data.get("status") == "success"
    app = create_data["data"]["app"]
    app_id = app["id"]
    # Confirm app isDisabled initially (isEnabled=False)
    assert app.get("isEnabled") is False

    try:
        # Step 3: Enable the disabled app via PATCH /api/v1/apps/enable/{appId}
        enable_url = f"{BASE_URL}/api/v1/apps/enable/{app_id}"
        patch_resp1 = requests.patch(enable_url, headers=headers, timeout=TIMEOUT)
        assert patch_resp1.status_code == 200, f"Enable app failed: {patch_resp1.text}"
        patch_data1 = patch_resp1.json()
        assert patch_data1.get("status") == "success" or "data" in patch_data1
        is_enabled1 = patch_data1.get("data", {}).get("isEnabled")
        assert is_enabled1 is True, f"Expected isEnabled=True after enabling, got {is_enabled1}"

        # Step 4: Test idempotency by enabling an already enabled app again
        patch_resp2 = requests.patch(enable_url, headers=headers, timeout=TIMEOUT)
        assert patch_resp2.status_code == 200, f"Idempotent enable app failed: {patch_resp2.text}"
        patch_data2 = patch_resp2.json()
        assert patch_data2.get("status") == "success" or "data" in patch_data2
        is_enabled2 = patch_data2.get("data", {}).get("isEnabled")
        assert is_enabled2 is True, f"Expected isEnabled=True on repeated enable, got {is_enabled2}"

    finally:
        # Cleanup: Delete the created app
        delete_url = f"{BASE_URL}/api/v1/apps/delete/{app_id}"
        delete_resp = requests.delete(delete_url, headers=headers, timeout=TIMEOUT)
        # Accept 200 or 204 or 404 if already deleted
        assert delete_resp.status_code in (200, 204, 404), f"Cleanup delete failed: {delete_resp.text}"

        # Optional: delete user (if API supported), skipping as no delete user endpoint described

test_patchappsenableapptogglesenabledstate()