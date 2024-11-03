FROM golang:1.23.2-alpine AS base
RUN apk --no-cache update

FROM base AS ci
WORKDIR /app/
COPY . .
RUN go mod tidy

FROM ci AS build
WORKDIR /app/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o entrypoint

FROM scratch
WORKDIR /
COPY --from=build /app/entrypoint .

ENTRYPOINT [ "/entrypoint" ]