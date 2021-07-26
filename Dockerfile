FROM golang:1.16

WORKDIR ./gogetent
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["gogetent"]