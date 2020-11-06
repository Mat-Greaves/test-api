FROM golang:alpine AS builder
WORKDIR $GOPATH/src/test-api/
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o /go/bin/api cmd/api/main.go

FROM scratch
COPY --from=builder /go/bin/api /go/bin/api
ENTRYPOINT [ "/go/bin/api" ]