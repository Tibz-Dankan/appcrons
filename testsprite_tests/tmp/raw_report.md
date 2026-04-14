
# TestSprite AI Testing Report(MCP)

---

## 1️⃣ Document Metadata
- **Project Name:** appcrons
- **Date:** 2026-04-14
- **Prepared by:** TestSprite AI Team

---

## 2️⃣ Requirement Validation Summary

#### Test TC001 postapipostcreateappwithvaliddata
- **Test Code:** [TC001_postapipostcreateappwithvaliddata.py](./TC001_postapipostcreateappwithvaliddata.py)
- **Test Error:** Traceback (most recent call last):
  File "/var/task/handler.py", line 258, in run_with_retry
    exec(code, exec_env)
  File "<string>", line 59, in <module>
  File "<string>", line 42, in test_post_api_post_create_app_with_valid_data
KeyError: 'data'

- **Test Visualization and Result:** https://www.testsprite.com/dashboard/mcp/tests/4e80a4b4-efca-4306-acc3-24efa395977c/70a6a650-9c94-40e1-aa53-c894f99161a1
- **Status:** ❌ Failed
- **Analysis / Findings:** {{TODO:AI_ANALYSIS}}.
---

#### Test TC002 patchappsupdateappwithvalidchanges
- **Test Code:** [TC002_patchappsupdateappwithvalidchanges.py](./TC002_patchappsupdateappwithvalidchanges.py)
- **Test Error:** Traceback (most recent call last):
  File "/var/task/handler.py", line 258, in run_with_retry
    exec(code, exec_env)
  File "<string>", line 68, in <module>
  File "<string>", line 22, in test_patch_apps_update_app_with_valid_changes
AssertionError: Missing 'user' in signup data: {'token': 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzYyMDQxMjQsImlhdCI6MTc3NjE3MTcyNCwidXNlcklkIjoiMGU0ZWVhMzMtZjNiNC00MTRiLThmZjMtZmNkYjExZmM3ZTk2In0.csgQY3tjqEI9ABV9jteLt3D4-UOnWUWObX_5obk4u8Q'}

- **Test Visualization and Result:** https://www.testsprite.com/dashboard/mcp/tests/4e80a4b4-efca-4306-acc3-24efa395977c/2cc13df6-1e50-4f80-9974-71d5f3926f3a
- **Status:** ❌ Failed
- **Analysis / Findings:** {{TODO:AI_ANALYSIS}}.
---

#### Test TC003 patchappsenableapptogglesenabledstate
- **Test Code:** [TC003_patchappsenableapptogglesenabledstate.py](./TC003_patchappsenableapptogglesenabledstate.py)
- **Test Error:** Traceback (most recent call last):
  File "/var/task/handler.py", line 258, in run_with_retry
    exec(code, exec_env)
  File "<string>", line 70, in <module>
  File "<string>", line 21, in test_patchappsenableapptogglesenabledstate
KeyError: 'user'

- **Test Visualization and Result:** https://www.testsprite.com/dashboard/mcp/tests/4e80a4b4-efca-4306-acc3-24efa395977c/e5f6a76f-da23-4ed9-ab0d-e4babf83b9ff
- **Status:** ❌ Failed
- **Analysis / Findings:** {{TODO:AI_ANALYSIS}}.
---

#### Test TC004 patchappsdisableapptogglesenabledstate
- **Test Code:** [TC004_patchappsdisableapptogglesenabledstate.py](./TC004_patchappsdisableapptogglesenabledstate.py)
- **Test Error:** Traceback (most recent call last):
  File "/var/task/handler.py", line 258, in run_with_retry
    exec(code, exec_env)
  File "<string>", line 67, in <module>
  File "<string>", line 36, in test_patch_apps_disable_app_toggles_enabled_state
KeyError: 'data'

- **Test Visualization and Result:** https://www.testsprite.com/dashboard/mcp/tests/4e80a4b4-efca-4306-acc3-24efa395977c/49e5ef39-9b9c-414a-964d-4cf5d0e9080f
- **Status:** ❌ Failed
- **Analysis / Findings:** {{TODO:AI_ANALYSIS}}.
---

#### Test TC005 deleteappsdeleteappwithvalidid
- **Test Code:** [TC005_deleteappsdeleteappwithvalidid.py](./TC005_deleteappsdeleteappwithvalidid.py)
- **Test Error:** Traceback (most recent call last):
  File "/var/task/handler.py", line 258, in run_with_retry
    exec(code, exec_env)
  File "<string>", line 65, in <module>
  File "<string>", line 23, in test_deleteappsdeleteappwithvalidid
AssertionError: User data missing in signup response

- **Test Visualization and Result:** https://www.testsprite.com/dashboard/mcp/tests/4e80a4b4-efca-4306-acc3-24efa395977c/a0bf1b01-8329-4112-a30d-0f47cdb5c29a
- **Status:** ❌ Failed
- **Analysis / Findings:** {{TODO:AI_ANALYSIS}}.
---

#### Test TC006 postrequestspostrequesttimewithvalidinterval
- **Test Code:** [TC006_postrequestspostrequesttimewithvalidinterval.py](./TC006_postrequestspostrequesttimewithvalidinterval.py)
- **Test Error:** Traceback (most recent call last):
  File "/var/task/handler.py", line 258, in run_with_retry
    exec(code, exec_env)
  File "<string>", line 68, in <module>
  File "<string>", line 33, in test_postrequestspostrequesttimewithvalidinterval
AssertionError: App creation failed: {"message":"App URL already exists!","status":"error"}


- **Test Visualization and Result:** https://www.testsprite.com/dashboard/mcp/tests/4e80a4b4-efca-4306-acc3-24efa395977c/9105adf4-cc12-4a46-8a65-494f121c1988
- **Status:** ❌ Failed
- **Analysis / Findings:** {{TODO:AI_ANALYSIS}}.
---

#### Test TC007 postrequestspostrequesttimewithinvalidinterval
- **Test Code:** [TC007_postrequestspostrequesttimewithinvalidinterval.py](./TC007_postrequestspostrequesttimewithinvalidinterval.py)
- **Test Error:** Traceback (most recent call last):
  File "/var/task/handler.py", line 258, in run_with_retry
    exec(code, exec_env)
  File "<string>", line 75, in <module>
  File "<string>", line 43, in test_post_requests_post_request_time_with_invalid_interval
AssertionError: App ID not returned

- **Test Visualization and Result:** https://www.testsprite.com/dashboard/mcp/tests/4e80a4b4-efca-4306-acc3-24efa395977c/f67a44f5-1575-41bd-bce8-784a85f828e4
- **Status:** ❌ Failed
- **Analysis / Findings:** {{TODO:AI_ANALYSIS}}.
---

#### Test TC008 getrequestsgetliveprovidesseestream
- **Test Code:** [TC008_getrequestsgetliveprovidesseestream.py](./TC008_getrequestsgetliveprovidesseestream.py)
- **Test Error:** Traceback (most recent call last):
  File "/var/task/handler.py", line 258, in run_with_retry
    exec(code, exec_env)
  File "<string>", line 75, in <module>
  File "<string>", line 62, in test_get_requests_get_live_provides_sse_stream
AssertionError

- **Test Visualization and Result:** https://www.testsprite.com/dashboard/mcp/tests/4e80a4b4-efca-4306-acc3-24efa395977c/1f28c64e-795c-43f4-b2b8-f374274e289e
- **Status:** ❌ Failed
- **Analysis / Findings:** {{TODO:AI_ANALYSIS}}.
---

#### Test TC009 postfeedbackpostfeedbackwithvaliddata
- **Test Code:** [TC009_postfeedbackpostfeedbackwithvaliddata.py](./TC009_postfeedbackpostfeedbackwithvaliddata.py)
- **Test Error:** Traceback (most recent call last):
  File "/var/task/handler.py", line 258, in run_with_retry
    exec(code, exec_env)
  File "<string>", line 49, in <module>
  File "<string>", line 23, in test_post_feedback_post_feedback_with_valid_data
AssertionError: Signup response missing 'user'

- **Test Visualization and Result:** https://www.testsprite.com/dashboard/mcp/tests/4e80a4b4-efca-4306-acc3-24efa395977c/5ef1016f-2faf-45e6-aa52-8f2a115428ef
- **Status:** ❌ Failed
- **Analysis / Findings:** {{TODO:AI_ANALYSIS}}.
---

#### Test TC010 getmetricsreturnsprometheusmetrics
- **Test Code:** [TC010_getmetricsreturnsprometheusmetrics.py](./TC010_getmetricsreturnsprometheusmetrics.py)
- **Test Visualization and Result:** https://www.testsprite.com/dashboard/mcp/tests/4e80a4b4-efca-4306-acc3-24efa395977c/7dde9c11-abd2-4d76-8da6-c93b8ad442e7
- **Status:** ✅ Passed
- **Analysis / Findings:** {{TODO:AI_ANALYSIS}}.
---


## 3️⃣ Coverage & Matching Metrics

- **10.00** of tests passed

| Requirement        | Total Tests | ✅ Passed | ❌ Failed  |
|--------------------|-------------|-----------|------------|
| ...                | ...         | ...       | ...        |
---


## 4️⃣ Key Gaps / Risks
{AI_GNERATED_KET_GAPS_AND_RISKS}
---