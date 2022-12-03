FROM golang:1.19

WORKDIR /wichteln

COPY . .

RUN go build

EXPOSE 4200

CMD [ "./wichteln" ]