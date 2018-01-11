package api

// Config etcd config
type Config struct {
	CAcert     string `form:"cacert"`
	ClientKey  string `form:"clientkey"`
	ClientCert string `form:"clientcert"`
	Endpoint   string `form:"endpoint" binding:"required"`
}
