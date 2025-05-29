FROM golang:1.24

# Set the working directory
WORKDIR /app

# Copy the api-backend directory
COPY . .

# Download the Go dependencies
RUN go mod download

# Build the API backend
RUN go build -o ./api ./main.go

# Expose the port
EXPOSE 80

# Run the API backend
CMD ["./app/api"]