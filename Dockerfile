FROM abilioesteves/gowebbuilder:1.2.0 AS BUILD

WORKDIR /app

ADD go.mod .
ADD go.sum .
ADD main.go .

RUN go mod download

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /generator .

FROM flaviostutz/etcd-registrar:1.0.1 AS REGISTRAR

FROM alpine:latest

EXPOSE 32856

COPY --from=BUILD /generator /generator
COPY --from=REGISTRAR /bin/etcd-registrar /bin/etcd-registrar

RUN chmod 777 /generator

ENV REGISTRY_ETCD_URL ""
ENV REGISTRY_ETCD_BASE "/services"
ENV REGISTRY_SERVICE "metrics-generator"
ENV REGISTRY_TTL "60"
ENV PORT "32865"

ADD startup.sh .
RUN chmod 777 startup.sh

CMD ["sh", "startup.sh"]
