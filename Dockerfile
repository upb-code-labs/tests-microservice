# -- Build --
FROM docker.io/library/golang:1.21.5-alpine3.19 AS builder

# Install upx
WORKDIR /source
RUN apk --no-cache add git upx

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN go build -o /source/bin/artifact

# -- Run --
FROM docker.io/library/alpine:3.18 AS runner

# Add non-root user
RUN adduser -D -h /opt/codelabs -s /sbin/nologin codelabs
WORKDIR /opt/codelabs
USER codelabs

# Copy binary and run
COPY --from=builder /source/bin/artifact /source/bin/artifact

# Run
EXPOSE 8080
ENTRYPOINT ["/source/bin/artifact"]