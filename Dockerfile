##
## Build Container
##
FROM golang:1.22-alpine as build


WORKDIR /tmp/fedi-games

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GOOS=linux go build -o ./out/fedi-games .


##
## Run Container
##
FROM alpine
RUN apk add ca-certificates

COPY --from=build /tmp/fedi-games/out/ /bin/

# This container exposes port 8080 to the outside world
EXPOSE 3000

WORKDIR /fedi-games

# Run the binary program produced by `go install`
ENTRYPOINT ["/bin/fedi-games"]