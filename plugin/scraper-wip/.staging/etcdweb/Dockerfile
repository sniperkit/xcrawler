FROM busybox

COPY etcdweb /root
COPY ui/ /root/ui
CMD chmod +x /root/etcdweb
EXPOSE 8080

WORKDIR /root
CMD ["./etcdweb"]