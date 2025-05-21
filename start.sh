> time.json
> time2.json

# if [ "$1" == "test" ]; then
    # go run ./cmd/mock
# elif  [ "$1" == "server" ]; then
    # go run ./cmd/server
# else
    # go run ./cmd/container-dsh
# fi

go run ./cmd/container-dsh --mode="$1"
jq . time.json > time2.json