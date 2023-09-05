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
	Cookie   *http.Cookie
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

	request := fmt.Sprintf("http://%s:%s/mxml?action=login&username=%s&secret=%s", ami.IP, ami.Port, ami.Username, ami.Secret)
	resp, err := http.Get(request)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()

	// odg, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(odg))

	ami.cookieReciever(resp)
}

func (ami *Amigo) WaitEvent() {

	client := &http.Client{}

	url := fmt.Sprintf("http://%s:%s/mxml?action=waitevent", ami.IP, ami.Port)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
	request.AddCookie(ami.Cookie)

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err.Error())
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(string(data))
}

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
