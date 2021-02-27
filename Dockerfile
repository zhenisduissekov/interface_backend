from ubuntu:20.04

COPY . /

EXPOSE 8088

VOLUME /data

CMD ./interface_backend
