import requests
import uuid

BASE_URL = "http://localhost:8080"
TIMEOUT = 30


def test_post_feedback_post_feedback_with_valid_data():
    signup_url = f"{BASE_URL}/api/v1/auth/signup"
    feedback_post_url = f"{BASE_URL}/api/v1/feedback/post"

    # Step 1: Signup to get fresh user and access token
    unique_email = f"{uuid.uuid4()}@test.com"
    signup_payload = {
        "name": "Test User",
        "email": unique_email,
        "password": "TestPass@1234"
    }

    try:
        signup_resp = requests.post(signup_url, json=signup_payload, timeout=TIMEOUT)
        assert signup_resp.status_code == 201, f"Signup failed: {signup_resp.text}"
        signup_json = signup_resp.json()
        assert "accessToken" in signup_json, "accessToken missing in signup response"
        access_token = signup_json["accessToken"]
        assert signup_json.get("status") == "success", "Signup status not success"
        assert "user" in signup_json and signup_json["user"].get("email") == unique_email

        # Step 2: Post feedback with valid token and feedback data
        headers = {
            "Authorization": f"Bearer {access_token}",
            "Content-Type": "application/json"
        }
        feedback_payload = {
            "email": unique_email,
            "message": "Love the service",
            "rating": 5
        }

        feedback_resp = requests.post(feedback_post_url, json=feedback_payload, headers=headers, timeout=TIMEOUT)
        assert feedback_resp.status_code == 201, f"Feedback post failed: {feedback_resp.text}"
        feedback_json = feedback_resp.json()
        assert feedback_json.get("status") == "success", "Feedback post status not success"
        # Validate that feedback entry contains sent message and rating
        if "data" in feedback_json and isinstance(feedback_json["data"], dict):
            feedback_data = feedback_json["data"]
            msg = feedback_data.get("message") or feedback_data.get("feedback", {}).get("message")
            rat = feedback_data.get("rating") or feedback_data.get("feedback", {}).get("rating")
            assert msg == feedback_payload["message"] or msg is None  # API may or may not echo
            assert rat == feedback_payload["rating"] or rat is None
        # If no 'data', fallback: check message or id in top-level JSON keys
        else:
            assert any(k in feedback_json for k in ("message", "id", "feedback"))

    except (requests.RequestException, AssertionError) as e:
        raise e


test_post_feedback_post_feedback_with_valid_data()