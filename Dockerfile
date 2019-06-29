FROM golang:1.12 as builder
RUN mkdir /src
WORKDIR /src
COPY go.mod .
RUN go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o yakitserver

FROM scratch
COPY --from=builder /src/yakitserver /bin/
CMD ["/bin/yakitserver"]
