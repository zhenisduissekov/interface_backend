from ubuntu:20.04

COPY . /

EXPOSE 8088

CMD ./interface_backend
