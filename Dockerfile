FROM golang:1.22

WORKDIR /bin

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api .

CMD ["/bin/api"]