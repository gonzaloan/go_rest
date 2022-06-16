ARG GO_VERSION=1.16.6

FROM golang:${GO_VERSION}-alpine AS builder

# Compile our files. without proxy
RUN go env -w GOPROXY=direct
RUN apk add --no-cache git # need to get dependencies
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /src
# copy go.mod and go.sum to container
COPY ./go.mod ./go.sum ./
RUN go mod download


COPY ./ ./

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /goapi

# this is the one to execute our app
FROM scratch AS runner
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY .env ./

COPY --from=builder /goapi /goapi

EXPOSE 5050

ENTRYPOINT ["/goapi"]


