wrk.method = "POST"
wrk.body   = '{"contract_id": "token","method": "invoke","function": "transfer","source_account": "account_2070","dest_account": "account_40791","amount": 1}'
wrk.headers["Content-Type"] = "application/json"
-- wrk.headers["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImVhZDI5YWU5YTQ3NjRlYjU5MjNiMjhmZTM3MjljMTZlIiwibmFtZSI6Inlhbmd6aGl5YW8iLCJqdGkiOiIxNjI2MjU0NTYxOTA5IiwiaWF0IjoxNjI2MjU0NTYxfQ.KJjhtUvAzT2ZV6Otai4ZLuNLdXMhOZEAse2s8N6prp8"


