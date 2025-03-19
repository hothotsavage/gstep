FROM alpine as builder
LABEL maintainer=wsdev
EXPOSE 8216
RUN echo "Asia/Shanghai" > /etc/timezone
COPY ./app_linux /app/app_linux
COPY ./config_release.json /app/config_release.json
# 修改文件的执行权限
RUN chmod +x /app/app_linux
WORKDIR /app
CMD ["/app/app_linux"]