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

// Amigo structure for holding all the neccessary data for backed to communicate with Asterisk and Frontend propperly
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

// NewManager intializes an Amigo structure
func NewManager() *Amigo {
	return &Amigo{
		Devices: make(map[string]string),
		Hub:     NewHub(),
	}
}

// SetConf sets up Asterisk login information, returns error if login fails
func (ami *Amigo) SetConf(ip, port, user, secret string) error {
	if user == "" || secret == "" {
		return errors.New("login credentials must be provided")
	}

	// if ip was not provided assume localhost
	if ip == "" {
		ip = "127.0.0.1"
	} else if net.ParseIP(ip) == nil { // if ip parser doesnt return a pointer to net.IP then an invalid adress string has been provided
		return errors.New("invalid ip address")
	}

	if port == "" {
		port = "5038" // default asterisk port for telnet connection
	}

	ami.IP, ami.Port, ami.Username, ami.Secret = ip, port, user, secret
	return nil
}

// Start tells amigo to log into the asterisk interface with perviously configured config through tcp socket. Returns a channel through which Amigo's EventListener will return unparsed events
func (ami *Amigo) Start() (chan []string, error) {

	err := ami.Initialize() // collect data from asterisk to initialize our backend's structure
	if err != nil {
		return nil, err
	}

	ami.Conn, err = net.Dial("tcp", ami.IP+":"+ami.Port) // establish websocket connection to asterisk
	if err != nil {
		return nil, err
	}

	// send a login with provided configuration, no failover
	loginCmd := fmt.Sprintf("Action: Login\r\nUsername: %s\r\nSecret: %s\r\n\r\n", ami.Username, ami.Secret)
	_, err = ami.Conn.Write([]byte(loginCmd))
	if err != nil {
		return nil, err
	}
	ami.Reader = bufio.NewReader(ami.Conn) // reader is put into structure since it will be used by another method

	// consume initial response (login response from asterisk), reduces complexity in later code
	for ami.Reader.Buffered() > 0 {
		response, _ := ami.Reader.ReadString('\n')
		fmt.Println(response)
	}

	ev := make(chan []string) // channel for spreading events through the program
	go ami.EventListener(ev)  // start listening for events

	return ev, nil
}

// EventListener purely listens for events and to forwards them to application for processing
func (ami *Amigo) EventListener(ev chan []string) {
	for {
		temp := make([]string, 0, 15) // buffer with predifined capacity to avoid reallocations
		for {
			response, _ := ami.Reader.ReadString('\n')
			temp = append(temp, response)
			if strings.TrimSpace(response) == "" {
				ev <- temp
				break
			}
		} // this loop splits response by line into a slice and sends it to a handler
	}
}

// EventHandler takes events from EventListener and calls appropriate functions from handler
func (ami *Amigo) EventHandler(resp []string) {
	event := trimInfo(resp[0]) // pull event name from data
	filter[event](resp, ami)   // call the coresponding function
}

// Function used for initialization of Amigo structure. Calls for asterisk data through http to simplify the process
func (ami *Amigo) Initialize() error {
	jar, _ := cookiejar.New(nil) // store cookie to call action after login without auth
	client := &http.Client{
		Jar: jar,
	}

	url := fmt.Sprintf("http://%s:8088/rawman?action=login&username=%s&secret=%s", ami.IP, ami.Username, ami.Secret)
	_, err := client.Get(url) // execute login
	if err != nil {
		return err
	}

	url = fmt.Sprintf("http://%s:8088/rawman?action=devicestatelist", ami.IP)
	response, err := client.Get(url) // get devices from asterisk
	if err != nil {
		return err
	}

	b, _ := io.ReadAll(response.Body)
	split := strings.Split(string(b), "\r\n") // split events
	ami.Active = 0

	for i := 5; i < len(split)-6; i += 4 {
		ami.Devices[trimInfo(split[i])] = trimInfo(split[i+1])
		if trimInfo(split[i+1]) != "UNAVAILABLE" {
			ami.Active++
		}
	} // add devices and their states to device list in amigo, count active devices

	url = fmt.Sprintf("http://%s:8088/rawman?action=bridgelist", ami.IP)
	response, err = client.Get(url) // get brdiges
	if err != nil {
		return err
	}
	b, _ = io.ReadAll(response.Body)

	split = strings.Split(string(b), "\r\n")
	ami.Bridges, _ = strconv.Atoi(trimInfo(split[len(split)-3])) // simply takes bridge count that asterisk provides in response

	return nil
}
