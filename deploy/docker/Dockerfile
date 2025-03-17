FROM alpine as builder
LABEL maintainer=wsdev
EXPOSE 8216
RUN echo "Asia/Shanghai" > /etc/timezone
COPY ./app_linux /app/app_linux
COPY ./config_release.json /app/config_release.json
COPY docker-start.sh /app/start.sh
WORKDIR /app
CMD ["/app/app_linux"]