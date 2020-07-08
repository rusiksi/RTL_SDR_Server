FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Artem Danilchenko <art_injener@mail.ru>"

WORKDIR /app

# Copy go mod and sum files
COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/RTL_SDR_Server ./cmd/RTL_SDR_TcpIpServer/main.go

EXPOSE 62000

# Command to run the executable
CMD ["bin/RTL_SDR_Server"]

# For instance, this command will give the ADD command access to anything in your current directory.
# docker build -f docker-files/Dockerfile .
# Update: Docker now allows having the Dockerfile outside the build context (fixed in 18.03.0-ce, https://github.com/docker/cli/pull/886). So you can also do something like
# docker build -f ../Dockerfile .