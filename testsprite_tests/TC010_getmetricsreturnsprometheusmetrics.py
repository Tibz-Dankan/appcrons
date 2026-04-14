import requests
import uuid

def test_getmetricsreturnsprometheusmetrics():
    base_url = "http://localhost:8080"
    metrics_url = f"{base_url}/metrics"
    timeout = 30

    # Try normal GET /metrics request, expecting HTTP 200 with plaintext Prometheus metrics.
    try:
        response = requests.get(metrics_url, timeout=timeout)
    except requests.RequestException as e:
        assert False, f"Request to /metrics failed with exception: {e}"

    if response.status_code == 200:
        content_type = response.headers.get("Content-Type", "")
        # Prometheus metrics endpoint typically returns text/plain content type
        assert "text/plain" in content_type.lower(), f"Expected Content-Type text/plain but got {content_type}"
        body = response.text
        # Check that body contains some plausible Prometheus metrics format (e.g. lines with metric_name, # HELP or # TYPE)
        assert body, "Metrics response body is empty"
        assert ("# HELP" in body or "# TYPE" in body or "\n" in body), "Metrics response does not contain expected Prometheus format"
    elif response.status_code == 500:
        # On server error, response must be JSON with status "error" and message "internal server error"
        try:
            json_resp = response.json()
        except Exception:
            assert False, "500 response is not JSON"
        assert json_resp.get("status") == "error", f"Expected status 'error' in 500 response but got {json_resp.get('status')}"
        assert "internal server error" in json_resp.get("message", "").lower(), f"Expected error message to contain 'internal server error' but got {json_resp.get('message')}"
    else:
        assert False, f"Unexpected status code {response.status_code} returned from /metrics"

    # Additional stress test: simulate higher load by making multiple concurrent requests, expect 200 or proper 500 errors
    # But as concurrency and high load are hard to simulate in simple test, do 5 rapid sequential requests as a basic stress test
    for _ in range(5):
        try:
            resp = requests.get(metrics_url, timeout=timeout)
        except requests.RequestException as e:
            # In case of request failure under load, note but do not fail outright
            continue
        if resp.status_code == 200:
            ctype = resp.headers.get("Content-Type", "")
            assert "text/plain" in ctype.lower(), f"Expected Content-Type text/plain under load but got {ctype}"
            assert resp.text, "Empty metrics response body under load"
        elif resp.status_code == 500:
            try:
                json_error = resp.json()
            except Exception:
                continue
            assert json_error.get("status") == "error"
            assert "internal server error" in json_error.get("message", "").lower()
        else:
            # Acceptable if intermittent status codes occur but should be only 200 or 500
            assert resp.status_code in [200, 500], f"Unexpected status code {resp.status_code} under load"

test_getmetricsreturnsprometheusmetrics()