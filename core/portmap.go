package core

import (
	"net"
	"fmt"
	"bufio"
	log "github.com/Sirupsen/logrus"
	"io"
	"strings"
	"strconv"
	"os"
	"errors"
	"sync"
	"../debug"
)

type ParseError struct {
	Source string
	Text   string
}

func (e *ParseError) Error() string {
	return e.Text + ": " + e.Source
}

type InteruptedReader struct {
	R    io.Reader        // underlying reader
	Done <-chan struct {} // max bytes remaining
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

}

type PortMap struct {
	Port        uint
	SourceIP    *net.IPNet
	Destination string
	ACL         ACLCheck

	locker      *sync.Mutex
	listener    net.Listener
	done        chan struct {}
}


func NewPortMap() *PortMap {
	portmap := &PortMap{}
	portmap.locker = &sync.Mutex{}
	return portmap
}

//func transferData (conn1, conn2 net.TCPConn) {

//}

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
		listener, err := net.Listen("tcp", ":" + strconv.Itoa(int(portmap.Port)))
		portmap.listener = listener
		portmap.done = make(chan struct {})
		var wg sync.WaitGroup
		//defer ln.Close()
		if err == nil {
			debug.Printf("net.Listen: %s", listener.Addr())
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
						log.Fatal(err)
					}
					//Handle connection
					debug.Printf("listener.Accept(): %s", conn.RemoteAddr()) // debug
					connOut, err := net.Dial("tcp", portmap.Destination)
					if err != nil {
						return
						log.Fatalln(err)
					}
					debug.Printf("net.Dial(): %s", connOut.RemoteAddr()) // debug
					go func(conn1, conn2 net.Conn) {
						defer conn.Close()
						defer connOut.Close()
						io.Copy(conn1, &InteruptedReader{conn2, portmap.done})
						debug.Printf("Close connections: %s, %s", conn1.RemoteAddr(), conn2.RemoteAddr())
					}(conn, connOut)
					go func(conn1, conn2 net.Conn) {
						defer conn.Close()
						defer connOut.Close()
						io.Copy(conn1, &InteruptedReader{conn2, portmap.done})
						debug.Printf("Close connections: %s, %s", conn1.RemoteAddr(), conn2.RemoteAddr())
					}(connOut, conn)
				}
			}()

			go func() {
				wg.Wait()
				wait.Done()
				portmap.locker.Lock()
				defer portmap.locker.Unlock()
				listener.Close()
				portmap.Done()
				debug.Print("listener.Close()") // debug
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

func (portmap *PortMap) InitFromStr(str string) (error) {
	fields := strings.Fields(str)
	return portmap.InitFromFields(fields)
}

func (portmap *PortMap) InitFromFields(fields []string) (error) {
	if len(fields) >= 4 {
		if !strings.EqualFold(fields[0], "portmap") {
			return &ParseError{fields[0], fmt.Sprintf("Tag \"%s\" don't vatid", fields[0])}
		}
		port, err := strconv.ParseUint(fields[1], 0, 16)
		if err != nil {
			return &ParseError{fields[1], err.Error()}
		}
		portmap.Port = PortNumber(port)
		_, sourseip, err := net.ParseCIDR(fields[2])
		if err != nil {
			return &ParseError{fields[2], err.Error()}
		}
		portmap.SourceIP = sourseip
		portmap.Destination = fields[3]
		if _, err := net.ResolveTCPAddr("tcp", fields[3]); err != nil {
			return &ParseError{fields[3], err.Error()}
		}
		return nil
	} else {
		return &ParseError{"", "Can't parse the string"}
	}
}

func (portmap *PortMap) Init(SourceIP string, Port uint, Destination string, ACL string) (error) {
	if !(Port >= 1 && Port <= 0xFFFF) {
		return errors.New("Incorrect number of port, must be in range 1-65 535: " + strconv.Itoa(Port))
	}
	portmap.Port = Port

	if strings.TrimSpace(SourceIP) != "" {
		_, sourceIP, err := net.ParseCIDR(SourceIP)
		if err != nil {
			return err.Error()
		}
		portmap.SourceIP = sourceIP
	} else {
		portmap.SourceIP = nil
	}

	portmap.Destination = Destination
	if _, err := net.ResolveTCPAddr("tcp", Destination); err != nil {
		return err.Error()
	}
	return nil
}

type PortMapList map[PortNumber]*PortMap

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

func (list PortMapList) InitFromConfig(config *Config) (error) {
	for _, item := range config.PortMap {
		portmap := NewPortMap()
		if err := portmap.Init(item.SourceIP, item.Port, item.Destination, item.ACL); err != {
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

