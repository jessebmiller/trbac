package proxy

import (
	"net/http/httputil"
	"github.com/jessebmiller/trbac/auth"
)

type trbacProxy struct {
	ListenPort int
	backendURL url.URL
	auth auth.Auth
}

func (proxy trbacProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if proxy.auth.May(r) {
		httputil.NewSingleHostReverseProxy(proxy.backendURL).ServeHTTP(w, r)
	} else {
		http.Error(w, "Not Authorized", 401)
	}
}

// ConfiguredProxy reads the config and returns  trbacProxy
func ConfiguredProxy() (proxyConfig, error) {

}

