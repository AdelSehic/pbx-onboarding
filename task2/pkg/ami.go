package amigo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
)

type Amigo struct {
	IP       string
	Port     string
	Username string
	Secret   string
	Conn     net.Conn
	Reader   *bufio.Reader
	Devices  map[string]string
	Active   int
	Bridges  int
	Hub      *WSHub
}

func NewManager() *Amigo {
	return &Amigo{
		Devices: make(map[string]string),
		Hub:     NewHub(),
	}
}
func (ami *Amigo) SetConf(ip, port, user, secret string) error {
	if user == "" || secret == "" {
		return errors.New("login credentials must be provided")
	}

	if ip == "" {
		ip = "127.0.0.1"
	} else if net.ParseIP(ip) == nil {
		return errors.New("invalid ip address")
	}

	if port == "" {
		port = "5038"
	}

	ami.IP, ami.Port, ami.Username, ami.Secret = ip, port, user, secret
	return nil
}

func (ami *Amigo) Start() (chan []string, error) {
	ami.Initialize()
	var err error

	ami.Conn, err = net.Dial("tcp", ami.IP+":"+ami.Port)
	if err != nil {
		return nil, err
	}

	loginCmd := fmt.Sprintf("Action: Login\r\nUsername: %s\r\nSecret: %s\r\n\r\n", ami.Username, ami.Secret)
	_, err = ami.Conn.Write([]byte(loginCmd))
	if err != nil {
		return nil, err
	}
	ami.Reader = bufio.NewReader(ami.Conn)

	// consume initial response
	for ami.Reader.Buffered() > 0 {
		response, _ := ami.Reader.ReadString('\n')
		fmt.Println(response)
	}

	ev := make(chan []string)
	go ami.EventListener(ev)

	return ev, nil
}

func (ami *Amigo) EventListener(ev chan []string) {
	for {
		temp := make([]string, 0, 15)
		for {
			response, _ := ami.Reader.ReadString('\n')
			temp = append(temp, response)
			if strings.TrimSpace(response) == "" {
				ev <- temp
				break
			}
		}
	}
}

func (ami *Amigo) EventHandler(resp []string) {
	event := trimInfo(resp[0])
	if function, ok := filter[event]; ok {
		function(resp, ami)
	}
}

func (ami *Amigo) Initialize() error {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	url := fmt.Sprintf("http://%s:8088/rawman?action=login&username=%s&secret=%s", ami.IP, ami.Username, ami.Secret)

	_, err := client.Get(url)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("http://%s:8088/rawman?action=devicestatelist", ami.IP)
	devs, err := client.Get(url)
	if err != nil {
		return err
	}
	b, _ := io.ReadAll(devs.Body)

	split := strings.Split(string(b), "\r\n")

	ami.Active = 0
	for i := 5; i < len(split)-6; i += 4 {
		ami.Devices[trimInfo(split[i])] = trimInfo(split[i+1])
		if trimInfo(split[i+1]) == "NOT_INUSE" {
			ami.Active++
			fmt.Println("Active: ", ami.Active)
		}
	}

	url = fmt.Sprintf("http://%s:8088/rawman?action=bridgelist", ami.IP)
	brs, err := client.Get(url)
	if err != nil {
		return err
	}
	b, _ = io.ReadAll(brs.Body)

	split = strings.Split(string(b), "\r\n")
	ami.Bridges, _ = strconv.Atoi(trimInfo(split[len(split)-3]))

	return nil
}
