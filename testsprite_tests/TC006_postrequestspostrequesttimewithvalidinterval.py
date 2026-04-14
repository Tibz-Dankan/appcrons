import requests
import uuid

BASE_URL = "http://localhost:8080"
TIMEOUT = 30

def test_postrequestspostrequesttimewithvalidinterval():
    # Signup a new user with a dynamic email to get token
    signup_url = f"{BASE_URL}/api/v1/auth/signup"
    user_email = f"user_{uuid.uuid4()}@example.com"
    signup_payload = {
        "name": "Test User",
        "email": user_email,
        "password": "Secure@1234"
    }
    signup_resp = requests.post(signup_url, json=signup_payload, timeout=TIMEOUT)
    assert signup_resp.status_code == 201, f"Signup failed: {signup_resp.text}"
    signup_data = signup_resp.json()
    assert signup_data["status"] == "success"
    user_token = signup_data["data"]["token"]
    auth_headers = {
        "Authorization": f"Bearer {user_token}"
    }

    # Create a new app with a unique name to get appId
    create_app_url = f"{BASE_URL}/api/v1/apps/post"
    app_payload = {
        "name": f"Test App for Request Time {uuid.uuid4()}",
        "url": "https://example.com/health",
        "requestInterval": "10"
    }
    create_app_resp = requests.post(create_app_url, json=app_payload, headers=auth_headers, timeout=TIMEOUT)
    assert create_app_resp.status_code == 201, f"App creation failed: {create_app_resp.text}"
    app_data = create_app_resp.json()
    assert app_data["status"] == "success"
    app = app_data["data"]["app"]
    assert app["isEnabled"] is False
    app_id = app["id"]

    # Define valid request time frame with startTime before endTime
    post_request_time_url = f"{BASE_URL}/api/v1/requests/post-request-time"
    rtf_payload = {
        "appId": app_id,
        "startTime": "06:00",
        "endTime": "23:00",
        "timezone": "Africa/Kampala"
    }

    try:
        post_rtf_resp = requests.post(post_request_time_url, json=rtf_payload, headers=auth_headers, timeout=TIMEOUT)
        assert post_rtf_resp.status_code == 201, f"Post request time failed: {post_rtf_resp.text}"
        rtf_data = post_rtf_resp.json()
        assert rtf_data["status"] == "success"
        rtf = rtf_data["data"]["requestTime"]
        assert rtf["appId"] == app_id
        assert rtf["startTime"] == "06:00"
        assert rtf["endTime"] == "23:00"
        assert rtf["timezone"] == "Africa/Kampala"
    finally:
        # Cleanup: delete the created app to avoid clutter
        delete_app_url = f"{BASE_URL}/api/v1/apps/delete/{app_id}"
        # Best effort delete, ignore exceptions for cleanup
        try:
            requests.delete(delete_app_url, headers=auth_headers, timeout=TIMEOUT)
        except Exception:
            pass

test_postrequestspostrequesttimewithvalidinterval()
