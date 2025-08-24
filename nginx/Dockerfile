FROM nginx:latest

RUN apt-get update -y \
  && apt-get install -y openssl \
  && apt-get clean

RUN mkdir -p /ssl \
  && openssl req -x509 -newkey rsa:4096 -sha256 -nodes -keyout /ssl/ssl.key -out /ssl/ssl.crt -subj "/CN=Memora" -days 2000

WORKDIR /etc/nginx

COPY nginx.conf ./
COPY default.conf ./conf.d/


