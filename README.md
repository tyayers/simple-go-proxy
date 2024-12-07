# Simple Go Proxy
This is based on various sample code to proxy traffic either in proxy mode (curl -x) or pass-through mode (forwarding request to host).

```sh
# start service
go run .

# proxy mode
curl -x http://localhost:8080 https://www.httpbin.org/anything
# returns anything output from httpbin.org

# forward mode
curl http://localhost:8080/anything -H "host: httpbin.org"
# returns anything output from httpbin.org
```