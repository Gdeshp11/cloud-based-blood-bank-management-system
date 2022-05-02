FROM golang:1.18-alpine AS build

WORKDIR /server
COPY go.mod ./
COPY go.sum ./
COPY main.go ./
COPY static ./static
RUN go mod download
RUN CGO_ENABLED=0 go build -o /bin/bbms_server

FROM scratch
COPY --from=build /bin/bbms_server /bbms_server
COPY --from=build /server/static  /static

ENTRYPOINT ["/bbms_server"]
