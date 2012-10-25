package main 

import ( 
        "bufio" 
        "log" 
        "net" 
		"net/http"
		"time"
		"sync/atomic"
		"flag"
		"fmt"
		"code.google.com/p/go.net/websocket"
)

//transaction counter
var counter uint64 = 0
//to count transactions per second
var diffTransactions uint64 = 0
var counterPartial uint64 = 0
var lastTps uint64 = 0
var stackLastTps = new(Stack)

//webserver port
var port *int = flag.Int("p", 31415, "Port to listen.")

//struct to share information with websocket clients
type T struct {
	Msg  string
	PosX int
	PosY int
}

func main() {	
	//parsing parameters
	flag.Parse()

	//websocket and http server
	go func() {
		//dispatch /json request to jsonServer handler
		http.Handle("/json", websocket.Handler(jsonServer))
		//dispath any other request to the fileserver handler
		http.Handle("/", http.FileServer(http.Dir("./html")))
		//listen websocket and http on port 31415
		err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
		if err != nil {
			log.Fatalf("listen error on http port: %v", err.Error())
		}
	}()
	
	//listen tcpip socket on port 5882
	ln, err := net.Listen("tcp", ":5882") 
	if err != nil { 
		log.Fatalf("listen error on tcpip port: %v", err) 
	}

	log.Printf("NyanCat-server Listening tpcip on port 5882")
	log.Printf("NyanCat-server Listening websocket and http on http://localhost:PI")
	log.Printf("NyanCat-server S2 U")
 
	//treat each connection received with ln.Accept()
	go func() {
		for {
			//try to accept a new socket connection 
			conn, err := ln.Accept() 
			if err != nil { 
				log.Fatalf("accept error: %v", err) 
			}
			go HandleConnection(conn) 
		}
	}()

	//count the number of transactions per second every second! DUH!
	go func(){
		for {
			//get the transaction counter
			counterPartial = atomic.LoadUint64(&counter)
			
			//how many transaction in the last second
			lastTps = counterPartial - diffTransactions
			stackLastTps.Push((int)(lastTps))
			
			//print to stdout
			log.Printf("NyanCat-server transaction-number:%v tps:%v\n", counterPartial, lastTps)
							
			//store the last counter to know how much new transactions were proccessed in the last second
			diffTransactions = counterPartial
				
			//sleep 1 second
			time.Sleep(time.Second * 1)
		}
	}()

	//eternal wait
	select {}
} 

//handle a tcp connection receive from the main loop
func HandleConnection(conn net.Conn) {
	//more one socket connected. Added to the counter
	atomic.AddUint64(&counter, 1) 
        
	//initialize a bufio to read the socket stream
	bufr := bufio.NewReader(conn) 
	for { 
		//wait to read the socket stream until \n
		line, err := bufr.ReadString('\n') 
		if err != nil { 
			return 
		}

		//sleep 2 second then echoe back the message
		time.Sleep(time.Second * 2)
		conn.Write([]byte(line))
				
		//sleep more 1 second and close the socket
		time.Sleep(time.Second * 1)
		conn.Close()
	}
}

//handle a websocket connection received from webserver
func jsonServer(ws *websocket.Conn) {
	fmt.Printf("jsonServer %#v\n", ws.Config())
	for {
		var msg T
		
		// Receive receives a text message serialized T as JSON.
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("recv:%#v\n", msg)
		
		// Send send a text message serialized T as JSON.
		msg.Msg = "pong"
		msg.PosY = 300
		//calculate the position of Nyan Cat 
		msg.PosX = stackLastTps.Pop().(int)
		
		// Send to the webclients
		err = websocket.JSON.Send(ws, msg)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("send:%#v\n", msg)
	}
}