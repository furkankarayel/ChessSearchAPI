# Stage 1: Build stage for Go project
FROM golang:1.18 as builder

# Set the working directory for Go build
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the Go source code
COPY . .

RUN go test -v ./...

# Change to the cmd directory and build the Go project
WORKDIR /app/cmd
RUN go build -o /parser-service



# Stage 2: Build stage for C++ program
FROM gcc:latest as cpp-builder

# Set the working directory for C++ build
WORKDIR /cpp-app

# Copy the C++ source code
COPY bin/scoutfish/src .

# Determine the architecture type and compile the C++ program
RUN  make build ARCH=x86-64

# Stage 3: Run stage
FROM gcc:latest



# Set the working directory for running the Go binary
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /parser-service /parser-service
COPY --from=builder /app/pgn /app/pgn
COPY --from=cpp-builder /cpp-app/scoutfish /app/scoutfish

EXPOSE 8080

# Command to run the Go binary
ENTRYPOINT ["/parser-service"]