package amigo

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"

	persJar "github.com/juju/persistent-cookiejar"
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
	jar, _ := persJar.New(&persJar.Options{Filename: "./authcookie"})
	return &Amigo{
		Client: &http.Client{
			Jar: jar,
		},
		Devices: make(map[string]string),
	}
}

func (ami *Amigo) action(action string) (string, error) {
	url := fmt.Sprintf("http://%s:%s/rawman?action=%s", ami.IP, ami.Port, action)

	resp, err := ami.Client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	ami.Client.Jar.(*persJar.Jar).Save()
	return string(data), nil
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
	reg := regexp.MustCompile(`Response: (\w+)`)
	defer ami.FetchDevices()

	resp, err := ami.action("login")
	if err != nil {
		return err
	}
	response := reg.FindStringSubmatch(resp)
	if response[1] == "Success" {
		fmt.Println("Authenticated using cookie, fetching devices...")
		return nil
	}

	resp, err = ami.action(fmt.Sprintf("%s&username=%s&secret=%s", "login", ami.Username, ami.Secret))
	if err != nil {
		return err
	}
	response = reg.FindStringSubmatch(resp)
	if response[1] == "Error" {
		return errors.New("couldn't authenticate")
	}

	fmt.Println("Authentication successful, fetching devices...")
	return nil
}

func (ami *Amigo) EventListener() (chan string, chan error) {
	ch := make(chan string)
	errchan := make(chan error)

	go func() {
		for {
			resp, err := ami.action("waitevent")
			if err != nil {
				errchan <- err
			}
			ch <- resp
		}
	}()

	return ch, errchan
}

func (ami *Amigo) EventHandler(resp string) {
	reg := regexp.MustCompile(`Event: (\w+)`)
	event := reg.FindStringSubmatch(resp)
	if function, ok := filter[event[1]]; ok {
		function(resp, ami)
	}
}

func (ami *Amigo) FetchDevices() error {
	reg := regexp.MustCompile(`: ([^\r\n]+)`)

	resp, err := ami.action("DeviceStateList")
	if err != nil {
		return err
	}
	events := strings.Split(resp, "\r\n\r\n")

	for i := 1; i < len(events)-2; i++ {
		vals := reg.FindAllStringSubmatch(events[i], -1)
		ami.Devices[vals[1][1]] = vals[2][1]
	}
	return nil
}

func (ami *Amigo) ListDevices() {
	for dev, state := range ami.Devices {
		fmt.Printf("dev: %s, state: %s\n", dev, state)
	}
}
