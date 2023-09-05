package amigo

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type Amigo struct {
	IP       string
	Port     string
	Username string
	Secret   string
}

func NewManager() *Amigo { return &Amigo{} }

func (ami *Amigo) SetConf(ip, port, user, secret string) error {

	if user == "" || secret == "" {
		return errors.New("login credentials must be provided")
	}

	if ip == "" {
		ip = "localhost"
	} else if net.ParseIP(ip) == nil {
		return errors.New("invalid ip address")
	}

	if port == "" {
		port = "8088"
	}

	ami.IP, ami.Port, ami.Username, ami.Secret = ip, port, user, secret
	return nil
}

func (ami *Amigo) Login() {

	f, _ := os.Create("cookie.txt")

	request := fmt.Sprintf("http://%s:%s/mxml?action=login&username=%s&secret=%s", ami.IP, ami.Port, ami.Username, ami.Secret)
	resp, err := http.Get(request)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()

	odg, _ := io.ReadAll(resp.Body)
	fmt.Println(string(odg))

	cookies := resp.Cookies()
	for _, ck := range cookies {
		ckv := fmt.Sprintf("%s:%s\tFALSE\t/\tFALSE\t%d\t%s\t\"%s\"\n", ami.IP, ami.Port, time.Now().Add(time.Duration(ck.MaxAge)).Unix(), ck.Name, ck.Value)
		f.Write([]byte(ckv))
	}
}

/*
		format := `
Name  string %s
Value string %s

Path       string %s
Domain     string %s
Expires    time.Time %s
RawExpires string %s

MaxAge   int %d
Secure   bool %v
HttpOnly bool %v\b
Raw      string %s
Unparsed []string %s
`
		fmt.Printf(format, ck.Name, ck.Value, ck.Path, ck.Domain, ck.Expires, ck.RawExpires, ck.MaxAge, ck.Secure, ck.HttpOnly, ck.Raw, ck.Unparsed)
*/
