FROM golang:bullseye

LABEL authors="arinazabrodina"

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /godocker

EXPOSE 8080

CMD ["/godocker"]