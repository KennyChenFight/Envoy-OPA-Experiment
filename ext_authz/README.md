# envoy-opa ext_authz example
+ reference: [envoy docs](https://www.envoyproxy.io/docs/envoy/latest/start/sandboxes/ext_authz)

## Scripts
### first-start
```.bash
docker-compose up --build -d
```
### refresh-start
```.bash
docker-compose down
docker-compose pull
docker-compose up --build -d
```

## Content
此範例實作了serviceA 往 serviceB 時由envoy Proxy擷取流量並問OPA是否給予通行，在給serviceB response。
詳細實作待寫QQ