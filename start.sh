> time.json
> time2.json

if [ "$1" == "test" ]; then
    go run ./cmd/mock
else
    go run ./cmd/container-dsh
fi
jq . time.json > time2.json