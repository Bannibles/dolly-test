# dolly-test

Test Web service built on Dolly

## Testing end-points

```.sh
curl --cacert etc/dev/certs/rootca/test_dolly_root_CA.pem https://localhost:8443/v1/teams
{"code":"unauthorized","message":"the \"guest\" role is not allowed"}
```

