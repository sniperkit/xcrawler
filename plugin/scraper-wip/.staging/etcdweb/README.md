## etcdweb

### etcdweb, show etcd data in web



#### Get code from github and run

`````
$ go get github.com/robinle/etcdweb
$ cd `go list -f '{{.Dir}}' github.com/robinle/etcdweb`
$ go run etcdweb.go -port 8080
# Open http://localhost:8080/
`````



#### Or run in docker container

```
$ docker pull robinle/etcdweb:v0.1
$ docker run -d -p 8080:8080 robinle/etcdweb:v0.1
# Open http://localhost:8080/ 
```



#### Open the web page

config the etcd endpoint


![](https://raw.githubusercontent.com/RobinLe/etcdweb/master/ui/pics/etcdweb-endpoint.png)

view all etcd key value from web

![](https://raw.githubusercontent.com/RobinLe/etcdweb/master/ui/pics/etcdweb-table.png)

