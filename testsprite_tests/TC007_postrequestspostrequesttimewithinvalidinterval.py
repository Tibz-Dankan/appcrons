import requests
import uuid

def test_post_requests_post_request_time_with_invalid_interval():
    base_url = "http://localhost:8080"
    signup_url = f"{base_url}/api/v1/auth/signup"
    signin_url = f"{base_url}/api/v1/auth/signin"
    post_request_time_url = f"{base_url}/api/v1/requests/post-request-time"
    apps_post_url = f"{base_url}/api/v1/apps/post"
    apps_delete_url_template = f"{base_url}/api/v1/apps/delete/{{}}"
    timeout = 30

    # Step 1: Signup dynamic user to avoid duplicate email
    unique_email = f"user_{uuid.uuid4()}@example.com"
    password = "Secure@1234"
    signup_payload = {
        "name": "Test User",
        "email": unique_email,
        "password": password
    }
    resp = requests.post(signup_url, json=signup_payload, timeout=timeout)
    assert resp.status_code == 201, f"Signup failed: {resp.status_code}, {resp.text}"
    token = resp.json().get("data", {}).get("token")
    assert token, "No token returned on signup"

    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }

    # Step 2: Create app to get valid appId
    # Use unique URL to avoid duplicate errors
    unique_url = f"https://example.com/{uuid.uuid4()}"
    # Change requestInterval to string due to API expectation
    app_payload = {
        "name": "Test App for Invalid Interval",
        "url": unique_url,
        "requestInterval": "10"
    }
    resp_app = requests.post(apps_post_url, headers=headers, json=app_payload, timeout=timeout)
    assert resp_app.status_code == 201, f"App creation failed: {resp_app.status_code}, {resp_app.text}"
    app_id = resp_app.json().get("data", {}).get("app", {}).get("id")
    assert app_id, "App ID not returned"

    try:
        # Step 3: Send POST /api/v1/requests/post-request-time with invalid interval (startTime >= endTime)
        invalid_intervals = [
            {"startTime": "06:00", "endTime": "06:00"},
            {"startTime": "23:00", "endTime": "06:00"},
            {"startTime": "12:30", "endTime": "12:00"}
        ]
        for interval in invalid_intervals:
            payload = {
                "appId": app_id,
                "startTime": interval["startTime"],
                "endTime": interval["endTime"],
                "timezone": "Africa/Kampala"
            }
            resp_post_time = requests.post(post_request_time_url, headers=headers, json=payload, timeout=timeout)
            assert resp_post_time.status_code == 400, (
                f"Expected 400 for invalid interval {interval}, got {resp_post_time.status_code}: {resp_post_time.text}"
            )
            json_resp = resp_post_time.json()
            assert json_resp.get("status") == "error", f"Expected error status, got {json_resp}"
            assert "startTime must be before endTime" in json_resp.get("message", ""), (
                f"Unexpected error message for interval {interval}: {json_resp.get('message')}"
            )
    finally:
        # Cleanup: delete created app
        del_url = apps_delete_url_template.format(app_id)
        del_resp = requests.delete(del_url, headers=headers, timeout=timeout)
        # Accept 200 or 204 success or 404 if already deleted
        assert del_resp.status_code in [200, 204, 404], f"Unexpected delete status {del_resp.status_code}"

test_post_requests_post_request_time_with_invalid_interval()
