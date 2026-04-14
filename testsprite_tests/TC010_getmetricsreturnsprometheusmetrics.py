import requests
import uuid
import time

BASE_URL = "http://localhost:8080"
TIMEOUT = 30

def test_get_metrics_returns_prometheus_metrics():
    # Perform GET /metrics under normal load
    url = f"{BASE_URL}/metrics"
    try:
        resp = requests.get(url, timeout=TIMEOUT)
    except requests.RequestException as e:
        assert False, f"Request to /metrics failed: {e}"
    assert resp.status_code in (200, 500), f"Unexpected status code {resp.status_code}"
    if resp.status_code == 200:
        # For 200 OK, content-type should be plaintext (Prometheus metrics)
        content_type = resp.headers.get("Content-Type","")
        # Typically prometheus metrics is text/plain; version=0.0.4 or similar
        # We allow any text content type starting with text/plain
        assert content_type.startswith("text/plain"), f"Expected Content-Type to start with 'text/plain', got {content_type}"
        # Body should contain common Prometheus metric format symbols like # HELP or # TYPE or metric names (rough check)
        body = resp.text
        assert body.strip() != "", "Response body is empty"
        plausible_prometheus_metric = any(keyword in body for keyword in ["# HELP", "# TYPE", "go_gc_duration_seconds", "process_cpu_seconds_total", "up"])
        assert plausible_prometheus_metric, "Response body does not look like Prometheus metrics"
    else:
        # 500 Internal Server Error: JSON body with status and message
        try:
            data = resp.json()
        except Exception:
            assert False, "Expected JSON body on 500 error"
        assert data.get("status") == "error", f"Expected status 'error' in JSON but got {data.get('status')}"
        assert "internal server error" in data.get("message","").lower(), f"Expected error message containing 'internal server error', got: {data.get('message')}"

    # Additionally, test high load by rapid repeated requests (e.g. 10 requests in quick succession)
    # It should still respond 200 or 500 with appropriate messages
    # We will attempt 10 requests and check responses carefully
    errors_encountered = 0
    for _ in range(10):
        try:
            resp = requests.get(url, timeout=TIMEOUT)
        except requests.RequestException as e:
            errors_encountered += 1
            continue
        if resp.status_code == 200:
            ct = resp.headers.get("Content-Type","")
            if not ct.startswith("text/plain"):
                errors_encountered += 1
            body = resp.text
            if body.strip() == "" or all(keyword not in body for keyword in ["# HELP", "# TYPE", "go_gc_duration_seconds", "process_cpu_seconds_total", "up"]):
                errors_encountered += 1
        elif resp.status_code == 500:
            try:
                data = resp.json()
                if data.get("status") != "error" or "internal server error" not in data.get("message","").lower():
                    errors_encountered += 1
            except Exception:
                errors_encountered += 1
        else:
            errors_encountered += 1
        # small delay to avoid flooding too fast
        time.sleep(0.1)
    # We allow some errors but should not have all errors failing
    assert errors_encountered < 10, f"All high load /metrics requests failed or were invalid: {errors_encountered}/10"

test_get_metrics_returns_prometheus_metrics()