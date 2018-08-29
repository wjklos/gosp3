FROM golang:latest 
RUN go get -v github.com/gin-gonic/gin
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
EXPOSE 7718 
RUN go build -o main . 
CMD ["/app/main"]