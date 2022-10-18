FROM golang as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

# Copy the go source
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o cherry-ipfs 

FROM gcr.io/distroless/static

WORKDIR /
COPY --from=builder /workspace/cherry-ipfs .
COPY --from=builder /workspace/.env ./.env

ENTRYPOINT ["/cherry-ipfs"]
