import requests
import uuid

BASE_URL = "http://localhost:8080"
SIGNUP_URL = f"{BASE_URL}/api/v1/auth/signup"
CREATE_APP_URL = f"{BASE_URL}/api/v1/apps/post"
DELETE_APP_URL_TEMPLATE = f"{BASE_URL}/api/v1/apps/delete/{{appId}}"

def test_post_api_post_create_app_with_valid_data():
    # Step 1: Signup to get fresh access token
    signup_email = f"{uuid.uuid4()}@test.com"
    signup_payload = {
        "name": "Test User",
        "email": signup_email,
        "password": "TestPass@1234"
    }
    try:
        signup_resp = requests.post(SIGNUP_URL, json=signup_payload, timeout=30)
        assert signup_resp.status_code == 201, f"Signup failed with status {signup_resp.status_code}"
        signup_data = signup_resp.json()
        assert "accessToken" in signup_data, "accessToken not in signup response"
        access_token = signup_data["accessToken"]
        assert signup_data.get("status") == "success", "Signup status not success"
    except (requests.RequestException, AssertionError) as e:
        raise e

    headers = {
        "Authorization": f"Bearer {access_token}",
        "Content-Type": "application/json"
    }

    # Step 2: Create a new app with valid data
    app_name = f"App {uuid.uuid4()}"
    app_url = f"https://app-{uuid.uuid4()}.onrender.com/active"
    create_app_payload = {
        "name": app_name,
        "url": app_url,
        "requestInterval": "10"
    }

    app_id = None
    try:
        create_resp = requests.post(CREATE_APP_URL, json=create_app_payload, headers=headers, timeout=30)
        assert create_resp.status_code == 201, f"Create app failed with status {create_resp.status_code}"
        create_data = create_resp.json()
        assert create_data.get("status") == "success", "Create app status not success"
        assert "data" in create_data and "app" in create_data["data"], "App data missing in response"
        app = create_data["data"]["app"]
        assert isinstance(app, dict), "App field is not a dict"
        assert app.get("name") == app_name, "App name in response does not match"
        assert app.get("url") == app_url, "App url in response does not match"
        assert app.get("requestInterval") == "10", "App requestInterval in response does not match"
        assert "isDisabled" in app, "isDisabled field missing in app response"
        # According to rules, newly created apps have isDisabled=true by default.
        assert app["isDisabled"] is True, "Newly created app isDisabled is not true by default"
        app_id = app.get("id")
        assert app_id is not None, "App ID missing in response"
    except (requests.RequestException, AssertionError) as e:
        raise e
    finally:
        # Cleanup: Delete the created app if app_id is known
        if app_id:
            try:
                del_resp = requests.delete(DELETE_APP_URL_TEMPLATE.format(appId=app_id), headers=headers, timeout=30)
                # Allow 200 success and 403 (if token mismatch, but should not happen here)
                if del_resp.status_code not in (200, 403, 404):
                    raise Exception(f"Unexpected status deleting app: {del_resp.status_code}")
            except requests.RequestException:
                pass  # Ignore exceptions in cleanup

test_post_api_post_create_app_with_valid_data()