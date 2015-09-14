//   Copyright (C) 2015 Evgeny M. Safonov
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package core

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

type ParseError struct {
	Source string
	Text   string
}

func (e *ParseError) Error() string {
	return e.Text + ": " + e.Source
}

type InteruptedReader struct {
	R    io.Reader       // underlying reader
	Done <-chan struct{} // max bytes remaining
}

func (r *InteruptedReader) Read(p []byte) (n int, err error) {
	select {
	case <-r.Done:
		return 0, io.EOF
	default:
		n, err = r.R.Read(p)
	}
	return
}

type PortNumber uint

type IPNet net.IPNet

func (ipnet *IPNet) UnmarshalJSON(b []byte) error {
	return nil
}

type PortMap struct {
	Port        int
	SourceIP    *net.IPNet
	Destination string
	ACL         ACLCheck

	locker   *sync.Mutex
	listener net.Listener
	done     chan struct{}
}

func NewPortMap() *PortMap {
	portmap := &PortMap{}
	portmap.locker = &sync.Mutex{}
	return portmap
}

//func transferData (conn1, conn2 net.TCPConn) {

//}

func (portmap *PortMap) String() string {
	sourceIP := "any"
	if portmap.SourceIP != nil {
		sourceIP = portmap.SourceIP.String()
	}
	return fmt.Sprintf("Port: %d; SourseIP: %s; Destination: %s", portmap.Port, sourceIP, portmap.Destination)
}

func (portmap *PortMap) Done() {
	portmap.locker.Lock()
	defer portmap.locker.Unlock()
	if portmap.done != nil {
		close(portmap.done)
		portmap.done = nil
	}
}

func (portmap *PortMap) Start(wait *sync.WaitGroup) {
	portmap.locker.Lock()
	defer portmap.locker.Unlock()

	if portmap.listener == nil {
		listener, err := net.Listen("tcp", ":"+strconv.Itoa(int(portmap.Port)))
		portmap.listener = listener
		portmap.done = make(chan struct{})
		var wg sync.WaitGroup
		//defer ln.Close()
		if err == nil {
			log.Debugf("net.Listen: %s", listener.Addr())
			wait.Add(1)
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer wait.Done()
				for {
					conn, err := listener.Accept()
					if err != nil {
						// handle error
						return
						log.Errorln(err)
					}
					//Handle connection
					log.Debugf("listener.Accept(): (Remote addr:%s; Local addr: %s)", conn.RemoteAddr(), conn.LocalAddr()) // debug
					acceptConnection := true
					if portmap.SourceIP != nil {
						//						if TCPConn, ok := conn.(*net.TCPConn); ok {
						if TCPAddr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
							acceptConnection = portmap.SourceIP.Contains(TCPAddr.IP)
							//log.Debugf("acceptConnection: %v; IP: %#v; Mask: %#v", acceptConnection, TCPAddr.IP, portmap.SourceIP)
						} else {
							acceptConnection = false
						}
					}
					if acceptConnection {
						connOut, err := net.Dial("tcp", portmap.Destination)
						if err != nil {
							return
							log.Errorln(err)
						}
						log.Debugf("net.Dial(): %s", connOut.RemoteAddr()) // debug
						go func(conn1, conn2 net.Conn) {
							defer conn.Close()
							defer connOut.Close()
							io.Copy(conn1, &InteruptedReader{conn2, portmap.done})
							log.Debugf("Close connections: %s, %s", conn1.RemoteAddr(), conn2.RemoteAddr())
						}(conn, connOut)
						go func(conn1, conn2 net.Conn) {
							defer conn.Close()
							defer connOut.Close()
							io.Copy(conn1, &InteruptedReader{conn2, portmap.done})
							//log.Debugf("Close connections: %s, %s", conn1.RemoteAddr(), conn2.RemoteAddr())
						}(connOut, conn)
					} else {
						conn.Close()
						log.Debugf("Reject the connection: %s, %s", conn.RemoteAddr(), conn.LocalAddr())
					}
				}
			}()

			go func() {
				wg.Wait()
				wait.Done()
				portmap.locker.Lock()
				defer portmap.locker.Unlock()
				listener.Close()
				portmap.Done()
				log.Debugln("listener.Close()") // debug
			}()

		} else {
			log.Fatal(err) // handle error
		}

	}
}

func (portmap *PortMap) Stop() {
	portmap.locker.Lock()
	defer portmap.locker.Unlock()

	portmap.Done()
	portmap.listener.Close()
	portmap.listener = nil
}

func (portmap *PortMap) Listener() net.Listener {
	return portmap.listener
}

//func (portmap *PortMap) InitFromStr(str string) (error) {
//	fields := strings.Fields(str)
//	return portmap.InitFromFields(fields)
//}

//func (portmap *PortMap) InitFromFields(fields []string) (error) {
//	if len(fields) >= 4 {
//		if !strings.EqualFold(fields[0], "portmap") {
//			return &ParseError{fields[0], fmt.Sprintf("Tag \"%s\" don't vatid", fields[0])}
//		}
//		port, err := strconv.ParseUint(fields[1], 0, 16)
//		if err != nil {
//			return &ParseError{fields[1], err.Error()}
//		}
//		portmap.Port = port
//		_, sourseip, err := net.ParseCIDR(fields[2])
//		if err != nil {
//			return &ParseError{fields[2], err.Error()}
//		}
//		portmap.SourceIP = sourseip
//		portmap.Destination = fields[3]
//		if _, err := net.ResolveTCPAddr("tcp", fields[3]); err != nil {
//			return &ParseError{fields[3], err.Error()}
//		}
//		return nil
//	} else {
//		return &ParseError{"", "Can't parse the string"}
//	}
//}

func (portmap *PortMap) Init(SourceIP string, Port int, Destination string, ACL string) error {
	if !(Port >= 1 && Port <= 0xFFFF) {
		return errors.New("Incorrect number of port, must be in range 1-65 535: " + strconv.Itoa(Port))
	}
	portmap.Port = Port

	trimmedSourceIP := strings.ToLower(strings.TrimSpace(SourceIP))
	portmap.SourceIP = nil
	if !(trimmedSourceIP == "" || trimmedSourceIP == "any") {
		_, sourceIP, err := net.ParseCIDR(SourceIP)
		if err != nil {
			return err
		}
		portmap.SourceIP = sourceIP
	}

	portmap.Destination = Destination
	if _, err := net.ResolveTCPAddr("tcp", Destination); err != nil {
		return err
	}
	return nil
}

type PortMapList map[int]*PortMap

func (list PortMapList) Add(portmap *PortMap) {
	if portmap != nil {
		list[portmap.Port] = portmap
	} else {
		log.Print("It is not possible to add an unspecified portmap(nil).")
	}
}

//func (list PortMapList) LoadLine(line string) (error) {
//	if string([]rune(line)[0]) != "#" { //line[0]
//		fields := strings.Fields(line)
//		if strings.EqualFold(fields[0], "portmap") {
//			portmap := NewPortMap()
//			if err := portmap.InitFromFields(fields); err == nil {
//				list.Add(portmap)
//			} else {
//				return err
//			}
//		}
//	}
//	return nil
//}

//func (list PortMapList) Load(filePath string) (error) {
//	debug.Println("Load config")
//	list.Stop()
//	for k := range list {
//		delete(list, k)
//	}
//
//	file, err := os.Open(filePath)
//	if err != nil {
//		return errors.New("Don't possible open settings file: " + filePath)
//	}
//	scan := bufio.NewScanner(file)
//	for scan.Scan() {
//		if err := list.LoadLine(scan.Text()); err != nil {
//			return err
//		}
//	}
//	return nil
//}

func (list PortMapList) InitFromConfig(config *Config) error {
	for _, item := range config.PortMap {
		portmap := NewPortMap()
		if err := portmap.Init(item.SourceIP, item.Port, item.Destination, item.ACL); err != nil {
			return err
		}
		list.Add(portmap)
	}
	return nil
}

func (list PortMapList) Start(wait *sync.WaitGroup) {
	for _, portmap := range list {
		portmap.Start(wait)
	}
}

func (list PortMapList) Stop() {
	for _, portmap := range list {
		portmap.Stop()
	}
}

func NewPortMapList() PortMapList {
	return PortMapList{}
}

//func NewPortMapListFromFile(filePath string) (list PortMapList, err error) {
//	list = NewPortMapList()
//	err = list.Load(filePath)
//	return
//}
