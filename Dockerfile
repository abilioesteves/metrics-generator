FROM abilioesteves/gowebbuilder:1.2.0 AS BUILD

WORKDIR /app

ADD go.mod .
ADD go.sum .
ADD main.go .

RUN go mod download

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /generator .

FROM alpine:latest

EXPOSE 32856

COPY --from=BUILD /generator /generator

RUN chmod 777 /generator

RUN echo "Starting the almighty Metrics Generator Tabajara..."

CMD [ "/generator" ]
