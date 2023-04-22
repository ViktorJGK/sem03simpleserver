package main

import (
	"fmt"
	"github.com/ViktorJGK/is105sem03/mycrypt"
	"github.com/ViktorJGK/minyr/yr"
	"io"
	"log"
	"net"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:8000")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}

					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))
					switch msg := string(dekryptertMelding); msg {
					//log.Println(msg)
					case "ping":
						kryptertMelding := mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03, 4)
						log.Println("Kryptert melding: ", string(kryptertMelding))
						_, err = conn.Write([]byte(string(kryptertMelding)))
					case "Kjevik":
						convertedLine, err := yr.CelsiusToFahrenheitLine(msg)
						if err != nil {
							log.Println("Error converting temperature:", err)
							return // or handle the error as appropriate for your use case
						}
						kryptertMelding := mycrypt.Krypter([]rune(fmt.Sprintf(convertedLine)), mycrypt.ALF_SEM03, 4)
						log.Println("Kryptert melding: ", string(kryptertMelding))
						_, err = conn.Write([]byte(string(kryptertMelding)))
						if err != nil {
							log.Println("Error writing to connection:", err)
							return
						}

					default:
						_, err = c.Write(buf[:n])
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
