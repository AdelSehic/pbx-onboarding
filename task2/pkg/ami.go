package amigo

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
)

type Amigo struct {
	IP       string
	Port     string
	Username string
	Secret   string
	Client   *http.Client
	Devices  map[string]string
}

func NewManager() *Amigo {
	jar, _ := cookiejar.New(nil)
	return &Amigo{
		Client: &http.Client{
			Jar: jar,
		},
		Devices: make(map[string]string),
	}
}

func (ami *Amigo) action(action string) string {
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
	request := fmt.Sprintf("%s&username=%s&secret=%s", ami.action("login"), ami.Username, ami.Secret)
	resp, err := ami.Client.Get(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	reg := regexp.MustCompile(`Response: (\w+)`)
	response := reg.FindStringSubmatch(string(data))
	if response[1] == "Error" {
		return errors.New("couldn't authenticate")
	}

	fmt.Println("Authentication successful, fetching devices...")
	ami.FetchDevices()
	return nil
}

func (ami *Amigo) EventListener() chan *http.Response {
	url := ami.action("waitevent")
	ch := make(chan *http.Response)

	go func() {
		for {
			resp, err := ami.Client.Get(url)
			if err != nil {
				log.Fatalf(err.Error())
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

func (ami *Amigo) FetchDevices() {
	resp, err := ami.Client.Get(ami.action("DeviceStateList"))
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	events := strings.Split(string(data), "\r\n\r\n")
	reg := regexp.MustCompile(`: ([^\r\n]+)`)
	for i := 1; i < len(events)-2; i++ {
		vals := reg.FindAllStringSubmatch(events[i], -1)
		ami.Devices[vals[1][1]] = vals[2][1]
	}
}

func (ami *Amigo) ListDevices() {
	for dev, state := range ami.Devices {
		fmt.Printf("dev: %s, state: %s\n", dev, state)
	}
}
