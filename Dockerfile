## Build
FROM golang:1.25-bookworm AS build

WORKDIR /app

COPY . ./

RUN go mod download
RUN make build

## Deploy
FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/configs /configs
COPY --from=build /app/output/server /server

EXPOSE 8082

ENTRYPOINT ["/server"]
