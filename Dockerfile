FROM go:1.12 as builder
RUN mkdir /src
WORKDIR /src
COPY . .
RUN go build -o yakitserver

FROM scratch
COPY --from=builder /src/yakitserver /bin/
CMD ["/bin/yakitserver"]
