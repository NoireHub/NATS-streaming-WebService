FROM golang:1.20-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -v ./cmd/webservice

CMD [ "./webservice" ]