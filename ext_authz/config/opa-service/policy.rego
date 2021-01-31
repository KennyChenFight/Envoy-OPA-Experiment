package envoy.authz

import input.attributes.request.http as http_request

default allow = false

allow = true {
  http_request.method == "GET"
  http_request.headers["x-current-user"] = "OPA"
}
