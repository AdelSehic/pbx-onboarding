package amigo

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

type Amigo struct {
	IP       string
	Port     string
	Username string
	Secret   string
	Conn     net.Conn
	Builder  strings.Builder
	Reader   *bufio.Reader
	Devices  map[string]string
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

func (ami *Amigo) Start() (chan string, error) {
	// defer ami.FetchDevices()
	// defer ami.FetchBridgeCount()
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

	ev := make(chan string)

	go ami.EventListener(ev)

	return ev, nil
}

func (ami *Amigo) EventListener(ev chan string) {
	for {
		for {
			response, _ := ami.Reader.ReadString('\n')
			ami.Builder.Write([]byte(response))
			if strings.TrimSpace(response) == "" {
				ev <- ami.Builder.String()
				ami.Builder.Reset()
				break
			}
		}
	}
}

func (ami *Amigo) EventHandler(resp string) {
	event := extractEventName(resp)
	fmt.Printf("Caught event: %s\n", event)
	if function, ok := filter[event]; ok {
		function(resp, ami)
	}
}

func extractEventName(eventData string) string {
	startIndex := strings.Index(eventData, "Event:")
	if startIndex == -1 {
		return ""
	}
	endIndex := strings.Index(eventData[startIndex:], "\r\n")
	if endIndex == -1 {
		return ""
	}
	eventName := eventData[startIndex+7 : startIndex+endIndex]
	return strings.TrimSpace(eventName)
}

func (ami *Amigo) Action(action string) (string, error) {
	request := fmt.Sprintf("Action: %s\r\n\r\n", action)

	_, err := ami.Conn.Write([]byte(request))
	if err != nil {
		return "", err
	}

	for {
		response, _ := ami.Reader.ReadString('\n')
		ami.Builder.WriteString(response)

		if strings.TrimSpace(response) == "" {
			temp := ami.Builder.String()
			ami.Builder.Reset()
			return temp, nil
		}
	}
}

// func (ami *Amigo) FetchDevices() error {
// 	reg := regexp.MustCompile(`: ([^\r\n]+)`)

// 	resp, err := ami.action("DeviceStateList")
// 	if err != nil {
// 		return err
// 	}
// 	events := strings.Split(resp, "\r\n\r\n")

// 	for i := 1; i < len(events)-2; i++ {
// 		vals := reg.FindAllStringSubmatch(events[i], -1)
// 		ami.Devices[vals[1][1]] = vals[2][1]
// 	}
// 	return nil
// }

// func (ami *Amigo) FetchBridgeCount() {
// 	reg := regexp.MustCompile("ListItems: ([^\r\n])+")
// 	resp, err := ami.action("bridgelist")
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	ami.Bridges, _ = strconv.Atoi(reg.FindStringSubmatch(resp)[1])
// }
