package amigo

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func (ami *Amigo) cookieReciever(resp *http.Response) {
	cookies := resp.Cookies()
	for _, ck := range cookies {
		ami.Cookie = ck
		ami.Cookie.Path, ami.Cookie.Domain = "/", "false"
		ami.saveCookie()
	}
}

func (ami *Amigo) saveCookie() {
	f, _ := os.Create("cookie.txt")
	defer f.Close()
	ckv := fmt.Sprintf("%s:%s\t%s\t%s\t%v\t%d\t%s\t\"%s\"\n", ami.IP, ami.Port, ami.Cookie.Domain, ami.Cookie.Path, ami.Cookie.Secure, time.Now().Add(time.Duration(ami.Cookie.MaxAge)).Unix(), ami.Cookie.Name, ami.Cookie.Value)
	f.Write([]byte(ckv))
}
