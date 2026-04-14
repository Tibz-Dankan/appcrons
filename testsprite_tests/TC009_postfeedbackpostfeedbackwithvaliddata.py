import requests
import uuid

BASE_URL = "http://localhost:8080"
TIMEOUT = 30

def test_post_feedback_post_feedback_with_valid_data():
    # Step 1: Sign up a new user to get token (dynamic email to avoid duplicates)
    signup_url = f"{BASE_URL}/api/v1/auth/signup"
    unique_email = f"user_{uuid.uuid4()}@example.com"
    signup_payload = {
        "name": "Test User",
        "email": unique_email,
        "password": "Secure@1234"
    }
    signup_resp = requests.post(signup_url, json=signup_payload, timeout=TIMEOUT)
    assert signup_resp.status_code == 201, f"Signup failed: {signup_resp.text}"
    signup_json = signup_resp.json()
    assert signup_json["status"] == "success"
    data = signup_json.get("data")
    assert data is not None, "Signup response missing 'data'"
    user = data.get("user")
    assert user is not None, "Signup response missing 'user'"
    user_id = user["id"]
    token = data.get("token")
    assert token is not None, "Signup response missing 'token'"
    headers = {"Authorization": f"Bearer {token}"}

    # Step 2: POST /api/v1/feedback/post with valid feedback data
    feedback_url = f"{BASE_URL}/api/v1/feedback/post"
    feedback_payload = {
        "message": "This is a test feedback message",
        "rating": 5
    }

    feedback_resp = requests.post(feedback_url, json=feedback_payload, headers=headers, timeout=TIMEOUT)
    assert feedback_resp.status_code == 201, f"Feedback creation failed: {feedback_resp.text}"
    feedback_json = feedback_resp.json()
    assert feedback_json["status"] == "success"
    data = feedback_json.get("data")
    assert data is not None, "Feedback response missing 'data'"
    assert "id" in data, "Feedback entry missing 'id'"
    assert data.get("message") == feedback_payload["message"]
    assert data.get("rating") == feedback_payload["rating"]

    # No cleanup since no delete feedback endpoint in PRD


test_post_feedback_post_feedback_with_valid_data()
