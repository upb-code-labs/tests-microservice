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

# Install OpenJDK 17
RUN apk --no-cache add openjdk17

# Install mvn
RUN apk --no-cache add maven

# Add non-root user
RUN addgroup -g 1000 codelabs
RUN adduser -D -h /opt/codelabs -u 1000 -G codelabs codelabs
WORKDIR /opt/codelabs
USER codelabs

# Copy binary and run
COPY --from=builder /source/bin/artifact /source/bin/artifact

# Create tests execution directory
RUN mkdir /opt/codelabs/tests-execution
ENV TESTS_EXECUTION_DIRECTORY /opt/codelabs/tests-execution

# Run
ENTRYPOINT ["/source/bin/artifact"]