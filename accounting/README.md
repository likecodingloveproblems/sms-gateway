# Acccounting Service

## Description
Other services like API Gateway can communicate with Accounting service with delivery layer that is choosed to be gRPC
To serve high throughput of requests, this service used a distributed cache layer like Redis
So in each seconds it must serve up to 1,500 requests
