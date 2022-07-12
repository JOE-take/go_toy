package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "9000"
	TYPE ="tcp"
)

func main(){
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil{
		log.Fatal(err)
		os.Exit(1)
	}

	defer listen.Close()
	conn, err := listen.Accept()
	fmt.Println("connected")
	if err != nil{
		log.Fatal()
		os.Exit(1)
	}
	endsignal := make(chan bool)
	go handler(conn, endsignal)
	<-endsignal
}

func handler(conn net.Conn, endsignal chan bool) {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan bool)
	go get_message(&conn, ch)
	go send_message(&conn, ch)
	<-ch
	fmt.Println("CONNECT CLOSED")
	conn.Close()
	endsignal <- true
}

func send_message(conn *net.Conn, ch chan bool){
	var text string
	for {
		fmt.Scan(&text)
		if text == "bye" {
			fmt.Println("EXIT SIGNAL DETECT FROM YOU")
			(*conn).Write([]byte("EXIT SIGNAL FROM OPPONENT" + "\n")) //終わることを相手に知らせる
			ch <- true
			break
		}
		(*conn).Write([]byte(text + "\n"))
	}
	(*conn).Write([]byte("CONNECT CLOSED" + "\n"))
	return 
}

func get_message(conn *net.Conn, ch chan bool){
	for{
		buffer := make([]byte, 1024)
		(*conn).Read(buffer)
		fmt.Printf(string(buffer))
		if string(string(buffer[0:3])) == "bye"{
			fmt.Println("EXIT SIGNAL FROM OPPONENT")
			(*conn).Write([]byte("EXIT SIGNAL FROM YOU" + "\n")) //終わることを相手に知らせる
			(*conn).Write([]byte("CONNECT CLOSED" + "\n")) //終わることを相手に知らせる
			ch <- true
			break
		}
	}
}