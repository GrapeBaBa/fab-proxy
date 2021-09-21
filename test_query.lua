wrk.method = "POST"
wrk.body   = '{"contract_id": "token","method": "invoke","function": "query","account": "account_40791"}'
wrk.headers["Content-Type"] = "application/json"
-- wrk.headers["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImVhZDI5YWU5YTQ3NjRlYjU5MjNiMjhmZTM3MjljMTZlIiwibmFtZSI6Inlhbmd6aGl5YW8iLCJqdGkiOiIxNjI2MjUyMjA1MDI1IiwiaWF0IjoxNjI2MjUyMjA1fQ.6c6KDRfO_A6MF_ZlZQfpZZqurB-nlNC-JIzWSHxqh5g"