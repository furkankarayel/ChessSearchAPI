# Stage 1: Build stage for Go project
FROM golang:1.18 AS go-builder

# Set the working directory for Go build
WORKDIR /app

# Copy the Go module files and download dependencies it's like package.json for Go
COPY go.mod ./
RUN go mod download

# Copy the Go source code (From the local computer into the stage working dir)
COPY . .

RUN go test -v ./...

# Change to the cmd directory (that's like a cd in terminal) 
# and build the Go project
WORKDIR /app/cmd

RUN go build -o /parser-service

# Stage 2: Get sources to build
FROM alpine:3.19.0 AS source-pull

# Set the working directory for C++ build
WORKDIR /cpp-app

# Install dependencies
RUN apk update && \
    apk add git 

# Pull pgn-extract and build binary
RUN git clone https://github.com/MichaelB7/pgn-extract.git 


# Pull scoutfish and build binary
RUN git clone https://github.com/mcostalba/scoutfish.git 

# Stage 3: Executable build stage
FROM gcc:14 AS cpp-builder

WORKDIR /cpp-app

COPY --from=source-pull /cpp-app/scoutfish /cpp-app/scoutfish
COPY --from=source-pull /cpp-app/pgn-extract /cpp-app/pgn-extract

# Pull pgn-extract and build binary
RUN cd pgn-extract/src && \
        make 


# Pull scoutfish and build binary
RUN cd scoutfish/src && \
        make build ARCH=x86-64


# Stage 3: Run stage
FROM gcc:14


# Set the working directory for running the Go binary
WORKDIR /app

# Copy all dependencies from the builder stage for the parser service
COPY --from=go-builder /parser-service /parser-service
COPY --from=go-builder /app/pgn /app/pgn

COPY --from=cpp-builder /cpp-app/scoutfish/src/scoutfish /app/scoutfish
COPY --from=cpp-builder /cpp-app/pgn-extract/src/pgn-extract /app/pgn-extract


EXPOSE 8080

# Command to run the Go binary
ENTRYPOINT ["/parser-service"]