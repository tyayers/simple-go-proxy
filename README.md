# Simple Go Proxy
This is based on various sample code to proxy traffic either in proxy mode (curl -x) or pass-through mode (forwarding request to host).

```sh
# start service
go run .

# test proxy mode
curl -x http://localhost:8080 https://www.httpbin.org/anything
# returns anything output from httpbin.org

# test forward mode
curl http://localhost:8080/anything -H "host: httpbin.org"
# returns anything output from httpbin.org

# test gcp vertex ai llm endpoint
curl -X POST "http://localhost:8080/v1beta1/projects/apigee-tlab1/locations/europe-west1/endpoints/openapi/chat/completions" \
     -H "Authorization: Bearer $(gcloud auth print-access-token)" \
     -H "Host: europe-west1-aiplatform.googleapis.com" \
     -H "Content-Type: application/json; charset=utf-8" \
	--data-binary @- << EOF

{
  "model": "google/gemini-1.5-flash-002",
  "stream": false,
  "messages": [{
    "role": "user",
    "content": "Write a story about a magic backpack."
  }]
}
EOF
```