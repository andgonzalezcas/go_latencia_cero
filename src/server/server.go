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

func main() {
	// 1. Se inicia el oyente
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error al iniciar el listener:", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Servidor escuchando en %s:%s...\n", CONN_HOST, CONN_PORT)
	defer listener.Close()

	// 2. Aqui llega la conexión que envia el cliente
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error al aceptar conexión:", err.Error())
		return
	}

	// esta conexión se mantiene con el cliente en todas las requests
	fmt.Println("Conexión persistente establecida. Iniciando bucle de latencia.")
	defer conn.Close()

	// 3. Bucle de Latencia Mínima (Read-Write Loop)
	buffer := make([]byte, MAX_MESSAGE_SIZE)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Conexión cerrada por el cliente. Terminando servidor.")
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
