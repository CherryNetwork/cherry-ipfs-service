FROM golang:1.16 as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

# Copy the go source
COPY main.go main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o cherry-ipfs main.go

FROM gcr.io/distroless/static

WORKDIR /
COPY --from=builder /workspace/cherry-ipfs .

ENTRYPOINT ["/cherry-ipfs"]