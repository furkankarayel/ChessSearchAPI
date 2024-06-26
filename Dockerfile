# Stage 1: Build stage for Go project
FROM golang:1.18 as go-builder

# Set the working directory for Go build
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the Go source code (From the local computer into the stage working dir)
COPY . .

RUN go test -v ./...

# Change to the cmd directory (that's like a cd in terminal) 
# and build the Go project
WORKDIR /app/cmd
RUN go build -o /parser-service



# Stage 2: Build stage for C++ program
FROM gcc:14 as cpp-builder

# Set the working directory for C++ build
WORKDIR /cpp-app

# Copy the C++ source code into /cpp-app
COPY bin/scoutfish/src .

# Determine the architecture type and compile the C++ program
RUN  make build ARCH=x86-64

# Stage 3: Run stage
FROM gcc:14



# Set the working directory for running the Go binary
WORKDIR /app

# Copy all dependencies from the builder stage for the parser service
COPY --from=go-builder /parser-service /parser-service
COPY --from=go-builder /app/pgn /app/pgn
COPY --from=cpp-builder /cpp-app/scoutfish /app/scoutfish

EXPOSE 8080

# Command to run the Go binary
ENTRYPOINT ["/parser-service"]