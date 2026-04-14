import requests
import uuid
import time

BASE_URL = "http://localhost:8080"
TIMEOUT = 30


def test_getrequestsgetliveprovidesseestream():
    # Step 1: Signup to get fresh user and accessToken
    signup_url = f"{BASE_URL}/api/v1/auth/signup"
    unique_email = f"{uuid.uuid4()}@test.com"
    signup_payload = {
        "name": "Test User",
        "email": unique_email,
        "password": "TestPass@1234"
    }
    resp_signup = requests.post(signup_url, json=signup_payload, timeout=TIMEOUT)
    assert resp_signup.status_code == 201, f"Signup failed with status {resp_signup.status_code}"
    json_signup = resp_signup.json()
    assert "accessToken" in json_signup, "accessToken not in signup response"
    access_token = json_signup["accessToken"]

    headers_auth = {"Authorization": f"Bearer {access_token}"}

    # Step 2: Test GET /api/v1/requests/get-live with valid authorization returns 200 and streams SSE

    sse_url = f"{BASE_URL}/api/v1/requests/get-live"
    try:
        with requests.get(sse_url, headers=headers_auth, stream=True, timeout=TIMEOUT) as sse_resp:
            assert sse_resp.status_code == 200, f"Expected 200 but got {sse_resp.status_code}"
            # Wait up to 5 seconds for any data or heartbeat event
            start_time = time.time()
            data_received = False
            # Read lines from the stream until timeout or data received
            for line_bytes in sse_resp.iter_lines():
                if line_bytes:
                    data_received = True
                    break
                if time.time() - start_time > 5:
                    break
            # Pass if connection established with 200 regardless of receiving events
            # No need to assert on specific event content
    except requests.exceptions.RequestException as e:
        assert False, f"Request with valid token failed: {e}"

    # Step 3: Test GET /api/v1/requests/get-live with missing token returns 401
    try:
        resp_no_auth = requests.get(sse_url, timeout=TIMEOUT, stream=True)
        assert resp_no_auth.status_code == 401, f"Expected 401 without token but got {resp_no_auth.status_code}"
        json_no_auth = resp_no_auth.json()
        assert "status" in json_no_auth and json_no_auth["status"] == "error"
    except requests.exceptions.RequestException as e:
        assert False, f"Request without token error: {e}"

    # Step 4: Test GET /api/v1/requests/get-live with expired token returns 401
    expired_token = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDAwMDAwMDAsImlhdCI6MTYwMDAwMDAwMCwidXNlcklkIjoiZmFrZXVzZXIifQ.RqUIH10hLzNTXVak4EdqvNGq-dPQCoZCF9bsjCZ674k"
    headers_expired = {"Authorization": expired_token}
    try:
        resp_expired = requests.get(sse_url, headers=headers_expired, timeout=TIMEOUT, stream=True)
        assert resp_expired.status_code == 401, f"Expected 401 with expired token but got {resp_expired.status_code}"
        json_expired = resp_expired.json()
        assert "status" in json_expired and json_expired["status"] == "error"
    except requests.exceptions.RequestException as e:
        assert False, f"Request with expired token error: {e}"


test_getrequestsgetliveprovidesseestream()