FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Artem Danilchenko <art_injener@mail.ru>"

WORKDIR /app

# Copy go mod and sum files
COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/ImitRTL_SDR_Reciver ./cmd/ImitRTL_SDR_Reciver/main.go ./cmd/ImitRTL_SDR_Reciver/imitObjectBuilder.go


# Command to run the executable
CMD ["bin/ImitRTL_SDR_Reciver"]
