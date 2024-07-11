FROM golang:latest
WORKDIR /APP
COPY . .
RUN go mod download
RUN go build -o bioskuy-image
ENTRYPOINT [ "/app/bioskuy-app" ]
