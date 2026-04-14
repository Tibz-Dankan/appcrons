import requests
import uuid

BASE_URL = "http://localhost:8080"
TIMEOUT = 30

def test_post_request_time_with_valid_interval():
    # Step 1: Signup - create a fresh user to get accessToken
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
    assert 'accessToken' in signup_json, f"No accessToken in signup response: {signup_resp.text}"
    access_token = signup_json['accessToken']

    headers = {
        "Authorization": f"Bearer {access_token}",
        "Content-Type": "application/json"
    }

    # Step 2: Create an app for the user
    create_app_url = f"{BASE_URL}/api/v1/apps/post"
    app_payload = {
        "name": f"App {uuid.uuid4()}",
        "url": f"https://app-{uuid.uuid4()}.onrender.com/active",
        "requestInterval": "10"
    }
    create_app_resp = requests.post(create_app_url, headers=headers, json=app_payload, timeout=TIMEOUT)
    assert create_app_resp.status_code == 201, f"App creation failed: {create_app_resp.text}"
    app_json = create_app_resp.json()
    assert app_json.get("status") == "success", f"App creation no success status: {create_app_resp.text}"
    app_data = app_json.get("data", {}).get("app")
    assert app_data and "id" in app_data, f"No app id in response: {create_app_resp.text}"
    app_id = app_data["id"]

    try:
        # Step 3: Post request time with valid interval
        post_request_time_url = f"{BASE_URL}/api/v1/requests/post-request-time"
        # Using exact HH:MM:SS format for start and end per rules
        request_time_payload = {
            "appId": app_id,
            "start": "06:00:00",
            "end": "23:00:00",
            "timeZone": "Africa/Kampala"
        }
        post_request_time_resp = requests.post(post_request_time_url, headers=headers, json=request_time_payload, timeout=TIMEOUT)
        assert post_request_time_resp.status_code == 201, f"Post request time failed: {post_request_time_resp.text}"
        response_json = post_request_time_resp.json()
        assert response_json.get("status") == "success", f"Request time creation no success status: {post_request_time_resp.text}"
        request_time_data = response_json.get("data", {}).get("requestTime")
        assert request_time_data is not None, f"No requestTime data in response: {post_request_time_resp.text}"
        # Assert required fields in the returned requestTime
        assert request_time_data.get("appId") == app_id, "Returned requestTime appId mismatch"
        assert request_time_data.get("start") == "06:00:00", "Returned requestTime start mismatch"
        assert request_time_data.get("end") == "23:00:00", "Returned requestTime end mismatch"
        assert request_time_data.get("timeZone") == "Africa/Kampala", "Returned requestTime timeZone mismatch"

    finally:
        # Cleanup: delete the created app
        delete_app_url = f"{BASE_URL}/api/v1/apps/delete/{app_id}"
        delete_resp = requests.delete(delete_app_url, headers=headers, timeout=TIMEOUT)
        # Deletion success is 200
        # If deletion fails, just log but do not raise to avoid masking test results
        if delete_resp.status_code != 200:
            print(f"Warning: Failed to delete app {app_id}: {delete_resp.status_code} {delete_resp.text}")

test_post_request_time_with_valid_interval()