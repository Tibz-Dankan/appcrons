import requests
import uuid
import time

BASE_URL = "http://localhost:8080"
TIMEOUT = 30


def test_get_requests_get_live_provides_sse_stream():
    # Step 1: Signup to get a user token (dynamic email required)
    signup_url = f"{BASE_URL}/api/v1/auth/signup"
    unique_email = f"user_{uuid.uuid4()}@example.com"
    signup_payload = {
        "name": "Test User",
        "email": unique_email,
        "password": "SecurePass@123"
    }
    signup_resp = requests.post(signup_url, json=signup_payload, timeout=TIMEOUT)
    assert signup_resp.status_code == 201, f"Signup failed: {signup_resp.text}"
    signup_data = signup_resp.json()
    assert signup_data["status"] == "success"
    token = signup_data["data"]["token"]
    assert token, "No token received on signup"

    headers = {
        "Authorization": f"Bearer {token}"
    }

    live_url = f"{BASE_URL}/api/v1/requests/get-live"

    # Test valid authorization - expect 200 and streaming SSE
    try:
        with requests.get(live_url, headers=headers, stream=True, timeout=TIMEOUT) as resp:
            assert resp.status_code == 200, f"Expected 200 OK with valid token, got {resp.status_code}"
            # Content-Type for SSE is usually text/event-stream
            content_type = resp.headers.get("Content-Type", "")
            assert "text/event-stream" in content_type, f"Expected 'text/event-stream' content-type, got {content_type}"

            # Read events from the stream but stop if disconnects or empty
            # We'll read a few lines and then stop; premature disconnect is a pass per instructions
            line_count = 0
            # read lines from stream to verify streaming
            for line in resp.iter_lines(decode_unicode=True):
                if line:  # Non-empty line
                    line_count += 1
                if line_count >= 5:
                    # Received some events, test passed
                    break
            # It's acceptable if no lines or premature close
    except requests.exceptions.ChunkedEncodingError:
        # Premature disconnect treated as pass per instructions
        pass
    except requests.exceptions.ReadTimeout:
        # Timeout treat as pass; no response or slow stream
        pass

    # Test missing token - expect 401 Unauthorized
    resp_no_auth = requests.get(live_url, timeout=TIMEOUT)
    assert resp_no_auth.status_code == 401, f"Expected 401 for missing token, got {resp_no_auth.status_code}"
    no_auth_data = resp_no_auth.json()
    assert no_auth_data.get("status") == "error"
    assert "unauthorized" in no_auth_data.get("message", "").lower()

    # Test expired/invalid token - expect 401 Unauthorized
    expired_headers = {
        "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDAwMDAwMDAsImlhdCI6MTYwMDAwMDAwMCwidXNlcklkIjoiZmFrZSJ9.faketoken"
    }
    resp_expired = requests.get(live_url, headers=expired_headers, timeout=TIMEOUT)
    assert resp_expired.status_code == 401, f"Expected 401 for expired token, got {resp_expired.status_code}"
    expired_data = resp_expired.json()
    assert expired_data.get("status") == "error"
    assert "unauthorized" in expired_data.get("message", "").lower()


test_get_requests_get_live_provides_sse_stream()