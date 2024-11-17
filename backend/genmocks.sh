set -e

# APIs

mockgen -package serviceMock \
-destination mocks/serviceMock/mocks.go \
-source service/service.go