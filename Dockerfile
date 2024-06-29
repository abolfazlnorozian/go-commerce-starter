FROM golang:1.18.10-alpine3.17
RUN apk update && \
    apk add --no-cache gcc musl-dev
    
WORKDIR /app

COPY go.mod go.sum ./
RUN GOPROXY=https://goproxy.cn go mod download

COPY . .
RUN go build -o main .

CMD ["./main"]