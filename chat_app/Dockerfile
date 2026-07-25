#  Use the official Go image
from golang:1.25.6

# set the working directory
WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# download dependencies
RUN go mod download

# copy the remaining project files
COPY . .

# Build the application
RUN go build -o server .

# Expose the application port
EXPOSE 8080

CMD [ "./server" ]

