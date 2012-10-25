package main 

import ( 
        "bufio" 
        "log" 
        "net" 
		"time"
		"sync/atomic"
)

var counter uint64 = 0
 
func main() {
        //listen tcp socket on port 5882
		ln, err := net.Listen("tcp", ":5882") 
        if err != nil { 
                log.Fatalf("listen error: %v", err) 
        }

 		log.Printf("NyanCat-server Listening on port 5882")
 
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

		go func(){
			var diffTransactions uint64 = 0
			for {
				//get the transaction counter
				counterPartial := atomic.LoadUint64(&counter)
				
				//print to stdout
				log.Printf("NyanCat-server transaction-number:%v tps:%v\n", counterPartial, counterPartial - diffTransactions)
				
				//store the last counter to know how much new transactions were proccessed in the last second
				diffTransactions = counterPartial
				
				//sleep 1 second
				time.Sleep(time.Second * 1)
			}
		}()

		//eternal wait
		select {}
} 

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