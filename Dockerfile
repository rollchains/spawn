# syntax=docker/dockerfile:1
#
# docker build . -t spawn:local
# docker run -it spawn:local

FROM golang:1.22.3 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Download local-ic (nested spawn add on)
RUN wget https://github.com/strangelove-ventures/interchaintest/releases/download/v8.8.0/local-ic && chmod +x local-ic
RUN mv ./local-ic /go/bin

# Build Spawn
RUN make build
RUN mv ./bin/spawn /go/bin

# Reduces the size of the final image from 4GB -> 0.25GB
FROM debian:12.6-slim as final

RUN apt update && apt install -y libc6-dev gcc make ca-certificates

# move spawn and local-ic to final
RUN mkdir -p /usr/local/bin
COPY --from=builder /go/bin/spawn /usr/local/bin/spawn
COPY --from=builder /go/bin/local-ic /usr/local/bin/local-ic

COPY --from=builder /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

CMD ["spawn"]