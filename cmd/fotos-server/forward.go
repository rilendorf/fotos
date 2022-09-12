package main

import (
	_ "embed"
	"encoding/binary"
	"fmt"
	"github.com/DerZombiiie/easyssh"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

var sshListener net.Listener

//go:embed id_rsa
var privateBytes []byte

const (
	//ip = "fotosgo.tk:22222"
	ip      = "localhost:22222"
	forward = "https://direct.fotosalbum.tk"

	minPort = 50000
	maxPort = 50050
)

var portsMu sync.RWMutex
var ports = make(map[int]*portData)

type portData struct {
	sync.RWMutex

	User, Pass string

	used bool
}

var usersMu sync.RWMutex
var users = make(map[string]string)

func getPort(user, pwd string) int {
	usersMu.Lock()
	users[user] = pwd
	usersMu.Unlock()

	portsMu.Lock()
	defer portsMu.Unlock()

	for i := minPort; i < maxPort; i++ {
		if _, ok := ports[i]; !ok {
			ports[i] = &portData{User: user, Pass: pwd}

			go func(i int) {
				time.Sleep(time.Second * 15)

				portsMu.Lock()
				defer portsMu.Unlock()
				p, ok := ports[i]
				if !ok {
					return
				}

				p.RLock()
				defer p.RUnlock()

				if !p.used {
					delete(ports, i)
				}
			}(i)

			return i
		}
	}

	log.Fatal("No free ports!")
	return -1
}

func main() {
	fmt.Printf("Started fotos forwarding service with portrange %d-%d\n", minPort, maxPort)
	fmt.Printf("configured SSH-address is %s\nand remote access url is %s\n", ip, forward)

	readUserDB()

	go startHttp()

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		fmt.Println("Error with Key", err)
	}

	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			fmt.Printf("connect of user %s with pass %s\n", c.User(), string(pass))

			usersMu.RLock()
			defer usersMu.RUnlock()

			if pwd, ok := users[c.User()]; ok && pwd == string(pass) {

				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %s", c.User())
		},
	}
	config.AddHostKey(private)
	easyssh.HandleChannel(easyssh.SessionRequest, easyssh.SessionHandler())
	easyssh.HandleChannel(easyssh.DirectForwardRequest, easyssh.DirectPortForwardHandler())
	easyssh.HandleRequestFunc(easyssh.RemoteForwardRequest, func(req *ssh.Request, sshConn ssh.Conn) {
		// check if allowed to do so
		fmt.Printf("User: %s is asking for a ssh tunnel\n", sshConn.User())

		t := easyssh.TcpipForward{}
		reply := (t.Port == 0) && req.WantReply
		ssh.Unmarshal(req.Payload, &t)
		addr := fmt.Sprintf("%s:%d", t.Host, t.Port)

		// check if allow listen:
		if t.Host != "127.0.0.1" && t.Host != "localhost" {
			return
		}

		portsMu.RLock()
		s, ok := ports[int(t.Port)]
		if !ok {
			portsMu.RUnlock()
			s.Unlock()
			return
		}

		s.Lock()

		if s.User != sshConn.User() {
			portsMu.RUnlock()
			s.Unlock()
			return
		}

		s.used = true
		s.Unlock()

		portsMu.RUnlock()

		ln, err := net.Listen("tcp", addr) //tie to the client connection

		if err != nil {
			fmt.Println("Unable to listen on address: ", addr)
			return
		}
		fmt.Println("Listening on address: ", ln.Addr().String())

		quit := make(chan bool)

		if reply { // Client sent port 0. let them know which port is actually being used

			_, port, err := easyssh.GetHostPortFromAddr(ln.Addr())
			if err != nil {
				return
			}

			b := make([]byte, 4)
			binary.BigEndian.PutUint32(b, uint32(port))
			t.Port = uint32(port)
			req.Reply(true, b)
		} else {
			req.Reply(true, nil)
		}

		go func() { // Handle incoming connections on this new listener
			for {
				select {
				case <-quit:

					return
				default:
					conn, err := ln.Accept()
					if err != nil { // Unable to accept new connection - listener likely closed
						continue
					}
					go func(conn net.Conn) {
						p := easyssh.DirectForward{}
						var err error

						var portnum int
						p.Host1 = t.Host
						p.Port1 = t.Port
						p.Host2, portnum, err = easyssh.GetHostPortFromAddr(conn.RemoteAddr())
						if err != nil {

							return
						}

						p.Port2 = uint32(portnum)
						ch, reqs, err := sshConn.OpenChannel(easyssh.ForwardedTCPReturnRequest, ssh.Marshal(p))
						if err != nil {
							fmt.Println("Open forwarded Channel: ", err.Error())
							return
						}
						go ssh.DiscardRequests(reqs)
						go func(ch ssh.Channel, conn net.Conn) {

							close := func() {
								ch.Close()
								conn.Close()

								// logger.Printf("forwarding closed")
							}

							go easyssh.CopyReadWriters(conn, ch, close)

						}(ch, conn)

					}(conn)
				}

			}

		}()
		sshConn.Wait()
		fmt.Println("Stop forwarding/listening on ", ln.Addr())
		ln.Close()
		quit <- true

		portsMu.Lock()
		defer portsMu.Unlock()

		delete(ports, int(t.Port))

	})
	//easyssh.ListenAndServe(":2022", config, nil)

	sshListener, err := net.Listen("tcp", ip)
	if err != nil {
		fmt.Println("SSH Error" + err.Error())
		return
	}
	easyssh.Serve(sshListener, config, nil)
}
