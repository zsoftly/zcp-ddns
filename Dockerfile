FROM golang:1.25-alpine AS build

ARG VERSION=dev

WORKDIR /src
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X github.com/zsoftly/zcp-ddns/internal/version.Version=${VERSION}" \
    -o /out/zcp-ddns ./cmd/zcp-ddns

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/zcp-ddns /zcp-ddns
USER nonroot:nonroot
ENTRYPOINT ["/zcp-ddns"]
