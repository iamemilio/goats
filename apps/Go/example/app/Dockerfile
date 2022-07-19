FROM golang

RUN mkdir /app
WORKDIR /app
COPY . .
RUN go build -o main .

CMD ["/app/main"]
