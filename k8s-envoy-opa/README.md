# K8S-Envoy-OPA
示意圖: ![arch](arch.png)
## Scripts
以MiniKube示範
```bash
minikube start
kubectl apply -f kubectl apply -f quick_start.yaml
kubectl expose deployment example-app --type=NodePort --name=example-app-service --port=8080
export SERVICE_PORT=$(kubectl get service example-app-service -o jsonpath='{.spec.ports[?(@.port==8080)].nodePort}')
export SERVICE_HOST=$(minikube ip)
export SERVICE_URL=$SERVICE_HOST:$SERVICE_PORT
echo $SERVICE_URL
export ALICE_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiZ3Vlc3QiLCJzdWIiOiJZV3hwWTJVPSIsIm5iZiI6MTUxNDg1MTEzOSwiZXhwIjoxNjQxMDgxNTM5fQ.K5DnnbbIOspRbpCr2IKXE9cPVatGOCBrBQobQmBmaeU"
export BOB_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJzdWIiOiJZbTlpIiwibmJmIjoxNTE0ODUxMTM5LCJleHAiOjE2NDEwODE1Mzl9.WCxNAveAVAdRCmkpIObOTaSd0AJRECY2Ch2Qdic3kU8"
```
Check that `Alice` can get employees **but cannot** create one.

```bash
curl -i -H "Authorization: Bearer "$ALICE_TOKEN"" http://$SERVICE_URL/people
curl -i -H "Authorization: Bearer "$ALICE_TOKEN"" -d '{"firstname":"Charlie", "lastname":"OPA"}' -H "Content-Type: application/json" -X POST http://$SERVICE_URL/people
```

Check that `Bob` can get employees and also create one.

```bash
curl -i -H "Authorization: Bearer "$BOB_TOKEN"" http://$SERVICE_URL/people
curl -i -H "Authorization: Bearer "$BOB_TOKEN"" -d '{"firstname":"Charlie", "lastname":"Opa"}' -H "Content-Type: application/json" -X POST http://$SERVICE_URL/people
```

Check that `Bob` **cannot** create an employee with the same firstname as himself.

```bash
curl -i  -H "Authorization: Bearer "$BOB_TOKEN"" -d '{"firstname":"Bob", "lastname":"Rego"}' -H "Content-Type: application/json" -X POST http://$SERVICE_URL/people
```
## 設計OPA bundle
**config.yaml**:

```yaml
services:
  - name: controller
    url: https://www.openpolicyagent.org
bundles:
  envoy/authz:
    service: controller
plugins:
  envoy_ext_authz_grpc:
    addr: :9191
    path: envoy/authz/allow
    dry-run: false
    enable-reflection: false
```

You can download the bundle and inspect it yourself:

```bash
mkdir example && cd example
curl -s -L https://www.openpolicyagent.org/bundles/envoy/authz | tar xzv
```

### Example Input

The `input` value defined for your policy will resemble the JSON below:

```json
{
  "attributes": {
    "source": {
      "address": {
        "socketAddress": {
          "address": "172.17.0.1",
          "portValue": 61402
        }
      }
    },
    "destination": {
      "address": {
        "socketAddress": {
          "address": "172.17.06",
          "portValue": 8000
        }
      }
    },
    "request": {
      "time": "2020-11-20T09:47:47.722473Z",
      "http": {
        "id":"13519049518330544501",
        "method": "POST",
        "headers": {
          ":authority":"192.168.99.206:30164",
          ":method":"POST",
          ":path":"/people?lang=en",
          "accept": "*/*",
          "authorization":"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJzdWIiOiJZbTlpIiwibmJmIjoxNTE0ODUxMTM5LCJleHAiOjE2NDEwODE1Mzl9.WCxNAveAVAdRCmkpIObOTaSd0AJRECY2Ch2Qdic3kU8",
          "content-length":"41",
          "content-type":"application/json",
          "user-agent":"curl/7.54.0",
          "x-forwarded-proto":"http",
          "x-request-id":"7bca5c86-bf55-432c-b212-8c0f1dc999ec"
        },
        "host":"192.168.99.206:30164",
        "path":"/people?lang=en",
        "protocol":"HTTP/1.1",
        "body":"{\"firstname\":\"Charlie\", \"lastname\":\"Opa\"}",
        "size":41
      }
    },
    "metadataContext": {}
  },
  "parsed_body":{"firstname": "Charlie", "lastname": "Opa"},
  "parsed_path":["people"],
  "parsed_query": {"lang": ["en"]},
  "truncated_body": false,
  "version": {
    "encoding":"protojson",
    "ext_authz":"v3"
  }
}
```
Note that this is the input [using the v3 API](#envoy-xds-v2-and-v2).

<details><summary>See here for an example of v2 input</summary>

```json
{
  "attributes":{
     "source":{
        "address":{
           "Address":{
              "SocketAddress":{
                 "PortSpecifier":{
                    "PortValue":61402
                 },
                 "address":"172.17.0.1"
              }
           }
        }
     },
     "destination":{
        "address":{
           "Address":{
              "SocketAddress":{
                 "PortSpecifier":{
                    "PortValue":8000
                 },
                 "address":"172.17.0.6"
              }
           }
        }
     },
     "request":{
        "http":{
           "id":"13519049518330544501",
           "method":"POST",
           "headers":{
              ":authority":"192.168.99.206:30164",
              ":method":"POST",
              ":path":"/people?lang=en",
              "accept":"*/*",
              "authorization":"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJzdWIiOiJZbTlpIiwibmJmIjoxNTE0ODUxMTM5LCJleHAiOjE2NDEwODE1Mzl9.WCxNAveAVAdRCmkpIObOTaSd0AJRECY2Ch2Qdic3kU8",
              "content-length":"41",
              "content-type":"application/json",
              "user-agent":"curl/7.54.0",
              "x-forwarded-proto":"http",
              "x-request-id":"7bca5c86-bf55-432c-b212-8c0f1dc999ec"
           },
           "host":"192.168.99.206:30164",
           "path":"/people?lang=en",
           "protocol":"HTTP/1.1",
           "body":"{\"firstname\":\"Charlie\", \"lastname\":\"Opa\"}",
           "size":41
        }
     }
  },
  "parsed_body":{"firstname": "Charlie", "lastname": "Opa"},
  "parsed_path":["people"],
  "parsed_query": {"lang": ["en"]},
  "truncated_body": false,
  "version": {
    "encoding":"encoding/json",
    "ext_authz":"v2"
  }
}
```
</details>

## 缺點
目前範例無法auto inject sidecar。可以自己建立一個webhook為K8S的MutatingWebhookConfiguration Object調用並動態塞入這些sidecar container

## Reference: https://github.com/open-policy-agent/opa-envoy-plugin

