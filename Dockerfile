# syntax=docker/dockerfile:1
#
# docker build . -t spawn:local
# docker run -it spawn:local

FROM golang:1.22.3 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Download and install local-ic to path
RUN wget https://github.com/strangelove-ventures/interchaintest/releases/download/v8.4.0/local-ic && chmod +x local-ic
RUN mv ./local-ic /go/bin

RUN make build
RUN mv ./bin/spawn /go/bin

# create a scratch image
FROM busybox:1.35.0 as final
RUN mkdir -p /usr/local/bin
COPY --from=builder /go/bin/spawn /usr/local/bin/spawn
COPY --from=builder /go/bin/local-ic /usr/local/bin/local-ic

# # run as busybopx
# RUN chmod +x /go/bin/spawn
# RUN chmod +x /go/bin/local-ic

CMD ["spawn"]