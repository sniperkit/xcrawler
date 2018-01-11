# Use goreman to run `go get github.com/mattn/goreman`

etcd1: etcd --name infra1 --listen-client-urls http://127.0.0.1:2379 --advertise-client-urls http://127.0.0.1:2379 --initial-cluster-token etcd-cluster-1 --initial-cluster 'infra1=http://127.0.0.1:12380' --initial-cluster-state new --enable-pprof

proxy: go run *.go ./shared/conf.d/providers.list.json