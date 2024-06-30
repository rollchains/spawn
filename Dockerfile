# syntax=docker/dockerfile:1
#
# docker build . -t spawn:local
# docker run -it spawn:local

FROM golang:1.22.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Download and install local-ic to path
RUN wget https://github.com/strangelove-ventures/interchaintest/releases/download/v8.4.0/local-ic && chmod +x local-ic
RUN mv ./local-ic /go/bin

RUN make build
RUN mv ./bin/spawn /go/bin

CMD ["spawn"]