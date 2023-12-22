# Start from the official Go image.
FROM golang:alpine

# Set the working directory inside the container.
WORKDIR /app

# Copy the Go files and other necessary files into the container.
COPY . .

# Build the application.
RUN go build -o main .

# Expose the port the app runs on.
EXPOSE 8080

# Command to run the executable.
CMD ["./main"]
