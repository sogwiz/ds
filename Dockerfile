FROM debian
RUN apt-get update
RUN apt-get install -y ca-certificates
COPY ./ds /ds
EXPOSE 3333
ENTRYPOINT ["/ds"]
