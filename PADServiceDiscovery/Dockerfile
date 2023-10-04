FROM golang:1.21-alpine

 # Set the working directory inside the container
WORKDIR /app

 # Copy the local package files to the container's workspace
COPY . .

 # Build the Go application inside the container
RUN go build -o main .

 # Expose port 5001
EXPOSE 5001

 # Run the Go application
 CMD ["./main"]