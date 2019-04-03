FROM harbor.haodai.net/base/alpine:3.7cgo
WORKDIR /app

MAINTAINER wenzhenglin(http://g.haodai.net/wenzhenglin/k8s-watcher.git)

COPY k8s-watcher /app

CMD /app/k8s-watcher
ENTRYPOINT ["./k8s-watcher"]

# EXPOSE 8080