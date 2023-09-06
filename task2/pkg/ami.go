package amigo

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
)

type Amigo struct {
	IP       string
	Port     string
	Username string
	Secret   string
	Cookie   *http.Cookie
	Client   *http.Client
	Devices  map[string]string
}

func NewManager() *Amigo {
	return &Amigo{
		Client:  &http.Client{},
		Devices: make(map[string]string),
	}
}

func (ami *Amigo) actionUrl(action string) string {
	return fmt.Sprintf("http://%s:%s/rawman?action=%s", ami.IP, ami.Port, action)
}

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

func (ami *Amigo) Login() error {

	request := fmt.Sprintf("%s&username=%s&secret=%s", ami.actionUrl("login"), ami.Username, ami.Secret)
	resp, err := http.Get(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ami.cookieReciever(resp)

	odg, _ := io.ReadAll(resp.Body)

	reg := regexp.MustCompile(`Response: (\w+)`)
	response := reg.FindStringSubmatch(string(odg))
	if response[1] == "Error" {
		return errors.New("couldn't authenticate")
	}

	fmt.Println("Authentication successful")
	return nil
}

func (ami *Amigo) EventListener() chan *http.Response {
	url := ami.actionUrl("waitevent")
	ch := make(chan *http.Response)

	go func() {
		for {
			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatalf(err.Error())
			}
			request.AddCookie(ami.Cookie)

			resp, err := ami.Client.Do(request)
			if err != nil {
				log.Fatal(err.Error())
			}
			ch <- resp
		}
	}()

	return ch
}

func (ami *Amigo) EventHandler(resp *http.Response) {

	reg := regexp.MustCompile(`Event: (\w+)`)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()

	event := reg.FindStringSubmatch(string(data))

	if function, ok := filter[event[1]]; ok {
		function(string(data), ami)
	}
}

// func (ami *Amigo) FetchDevices() {
// 	url := ami.actionUrl("DeviceStateList")

// }

func (ami *Amigo) ListDevices() {
	for dev, state := range ami.Devices {
		fmt.Printf("dev: %s, state: %s\n", dev, state)
	}
}
