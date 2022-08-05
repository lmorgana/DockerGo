FROM golang:1.18

WORKDIR /app

COPY ./htmlTemplate.go .
COPY ./html.gohtml .

COPY go.mod ./
RUN go mod download github.com/lib/pq

RUN go build htmlTemplate.go
CMD ./htmlTemplate postgres postgres html.gohtml