FROM alpine:3.4

RUN mkdir discovery

COPY . discovery

WORKDIR "discovery"

RUN chmod +x discovery

ENTRYPOINT ["./discovery"]

CMD ["serve"]