package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/lcpu-club/kube-dify-gui/internal/client"
)

func main() {
	c, err := client.NewHPCGame("sk:560a8c68-cda1-4783-a293-1378f427cce3:7O367iXWaiYY_DYZc4NaOnpEqmRploaOEDnv")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Do Port Forward")
	_, o, e, _ := c.DoPortForward("localhost", "28880")
	io.Copy(os.Stdout, o)
	io.Copy(os.Stderr, e)
}
