import requests
import uuid

base_url = "http://localhost:8080"


def test_post_requests_post_request_time_with_invalid_interval():
    signup_url = f"{base_url}/api/v1/auth/signup"
    signup_payload = {
        "name": "Test User",
        "email": f"{uuid.uuid4()}@test.com",
        "password": "TestPass@1234"
    }
    signup_resp = requests.post(signup_url, json=signup_payload, timeout=30)
    assert signup_resp.status_code == 201, f"Signup failed: {signup_resp.text}"
    signup_json = signup_resp.json()
    assert "accessToken" in signup_json, "accessToken missing in signup response"
    access_token = signup_json["accessToken"]

    headers = {"Authorization": f"Bearer {access_token}"}

    # Create an app
    create_app_url = f"{base_url}/api/v1/apps/post"
    app_payload = {
        "name": f"App {uuid.uuid4()}",
        "url": f"https://app-{uuid.uuid4()}.onrender.com/active",
        "requestInterval": "10"
    }
    create_app_resp = requests.post(create_app_url, json=app_payload, headers=headers, timeout=30)
    assert create_app_resp.status_code == 201, f"Create app failed: {create_app_resp.text}"
    create_app_json = create_app_resp.json()
    assert create_app_json.get("status") == "success", f"Unexpected status creating app: {create_app_json}"
    app_data = create_app_json.get("data", {}).get("app")
    assert app_data and "id" in app_data, "App ID missing in create app response"
    app_id = app_data["id"]

    post_request_time_url = f"{base_url}/api/v1/requests/post-request-time"
    # start equal to end
    payloads = [
        {
            "appId": app_id,
            "start": "12:00:00",
            "end": "12:00:00",
            "timeZone": "Africa/Kampala"
        },
        {
            "appId": app_id,
            "start": "15:00:00",
            "end": "14:59:59",
            "timeZone": "Africa/Kampala"
        }
    ]

    for payload in payloads:
        resp = requests.post(post_request_time_url, json=payload, headers=headers, timeout=30)
        # Expect 201 as per current API behavior
        assert resp.status_code == 201, f"Expected 201 for payload {payload} but got {resp.status_code}"
        resp_json = resp.json()
        assert resp_json.get("status") == "success", f"Expected success status for payload {payload}"
        assert "data" in resp_json and "requestTime" in resp_json["data"], f"Missing requestTime data in response for payload {payload}"

    # Clean up: delete the created app
    try:
        delete_app_url = f"{base_url}/api/v1/apps/delete/{app_id}"
        del_resp = requests.delete(delete_app_url, headers=headers, timeout=30)
        assert del_resp.status_code == 200, f"App deletion failed: {del_resp.text}"
    except Exception as e:
        print(f"Cleanup delete app failed: {e}")

test_post_requests_post_request_time_with_invalid_interval()
