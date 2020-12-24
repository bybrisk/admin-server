FROM golang:latest
LABEL maintainer="Shashank P. Sharma <shashanksharma191098@gmail.com>"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o admin
EXPOSE 80
CMD ["./admin"]