FROM golang:1.23
COPY . .
RUN go build -o server .
CMD [ "./server" ]
