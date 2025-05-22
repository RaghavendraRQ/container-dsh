> time.json
> time2.json

# if [ "$1" == "test" ]; then
    # go run ./cmd/mock
# elif  [ "$1" == "server" ]; then
    # go run ./cmd/server
# else
    # go run ./cmd/container-dsh
# fi

export PORT=":8080"
export CLIENT_URL="http://localhost:3000" # Change this later...

go run ./cmd/container-dsh --mode="$1"
jq . time.json > time2.json