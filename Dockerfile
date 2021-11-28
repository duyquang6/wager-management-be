FROM golang:1.17-alpine AS build-env
RUN apk --no-cache add build-base git mercurial gcc bash
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build
RUN make build.migrate

FROM alpine
WORKDIR /app
COPY --from=build-env /app/bin/app /app/
COPY --from=build-env /app/bin/migrate /app/
ENTRYPOINT ./app
