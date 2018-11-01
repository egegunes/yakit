FROM golang:1.11 as builder
RUN mkdir /src
WORKDIR /src
ADD . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/yakitcmd .
RUN chmod +x /src/yakitcmd

FROM scratch
COPY --from=builder /src/yakitcmd .
CMD ["./yakitcmd"]