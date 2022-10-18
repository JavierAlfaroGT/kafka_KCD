FROM golang
WORKDIR /
COPY . .
RUN go mod download
EXPOSE 7000
CMD ["go", "run", "consumer.go"]