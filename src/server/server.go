package main

import (
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST        = "localhost"
	CONN_PORT        = "8080"
	CONN_TYPE        = "tcp"
	RESPONSE         = "response"
	MAX_MESSAGE_SIZE = 64
)

func createConnection(host, port, connType string) net.Listener {
	listener, err := net.Listen(connType, host+":"+port)
	if err != nil {
		fmt.Println("Error al iniciar el listener:", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Servidor escuchando en %s:%s...\n", CONN_HOST, CONN_PORT)
	return listener
}

func handleRequest(conn net.Conn) {
	remoteAddr := conn.RemoteAddr()
	fmt.Printf("Conexión persistente establecida con: %s\n", remoteAddr)
	defer conn.Close()

	// 3. Bucle de Latencia Mínima (Read-Write Loop)
	buffer := make([]byte, MAX_MESSAGE_SIZE)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Conexión cerrada con %s\n", remoteAddr)
			return
		}

		if n > 0 {
			_, err := conn.Write([]byte(RESPONSE))
			if err != nil {
				fmt.Println("Error al enviar respuesta:", err.Error())
				return
			}

			receivedMessage := string(buffer[:n])
			fmt.Printf("Estímulo recibido: %s\n", receivedMessage)
		}
	}
}

func main() {
	// 1. Se inicia el oyente
	listener := createConnection(CONN_HOST, CONN_PORT, CONN_TYPE)
	defer listener.Close()

	// 2. Aqui llega la conexión que envia el cliente
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error al aceptar conexión:", err.Error())
		}

		// ahora tenemos un manejo de multiples clientes de ser requerido.
		go handleRequest(conn)
	}
}
