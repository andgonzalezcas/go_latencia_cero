package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
	CONN_TYPE = "tcp"

	MESSAGE          = "ping"
	MAX_MESSAGE_SIZE = 64

	NUM_REQUESTS = 10
	LOG_FILE     = "latency_log.txt"
)

func createConnection(host, port, connType string) net.Conn {
	addr := host + ":" + port
	conn, err := net.Dial(connType, addr)
	if err != nil {
		fmt.Println("Error al conectar con el servidor en", addr, ":", err.Error())
		os.Exit(1)
	}
	fmt.Println("Conexión establecida con", addr)
	return conn
}

func createFileConnection(file string, init_time string) *os.File {
	logFile, err := os.Create(file + init_time)
	if err != nil {
		fmt.Println("Error al crear el archivo de log:", err)
		os.Exit(1)
	}

	fmt.Fprintf(logFile, "--- Inicio de prueba: %s ---\n", init_time)
	return logFile
}

func runLatencyTest(conn net.Conn, numRequests int, logFile *os.File) time.Duration {
	var totalLatency time.Duration
	buffer := make([]byte, MAX_MESSAGE_SIZE)

	for i := 1; i <= NUM_REQUESTS; i++ {
		messageBytes := fmt.Appendf(nil, "%s %d", MESSAGE, i)
		start := time.Now()

		_, err := conn.Write(messageBytes)
		if err != nil {
			fmt.Printf("Error al enviar estímulo en iteración %d: %s.\n", i, err.Error())
			break
		}

		_, err = conn.Read(buffer)
		if err != nil {
			fmt.Printf("Error al leer request en iteración %d: %s.\n", i, err.Error())
			break
		}

		end := time.Now()
		latency := end.Sub(start) // end - start = tiempo de latencia
		totalLatency += latency

		status := "OK"
		latencyMs := float64(latency.Microseconds()) / 1000.0
		fmt.Fprintf(logFile, "%-9d | %-13.3f | %s\n", i, latencyMs, status)

		if i%100 == 0 || i == 1 {
			fmt.Printf("Iteración %d: Latencia: %.3f ms. Status: %s\n", i, latencyMs, status)
		}
	}

	return totalLatency
}

func createFinalReport(logFile *os.File, totalLatency time.Duration) {
	avgLatency := totalLatency / time.Duration(NUM_REQUESTS)
	avgLatencyMs := float64(avgLatency.Microseconds()) / 1000.0

	// Reporte final
	fmt.Println("\n--- PRUEBA FINALIZADA ---")
	fmt.Printf("Total de peticiones: %d\n", NUM_REQUESTS)
	fmt.Printf("Latencia Promedio: %.3f ms\n", avgLatencyMs)

	fmt.Fprintf(logFile, "--- Latencia Promedio: %.3f ms ---\n", avgLatencyMs)
	fmt.Printf("Revisa el archivo '%s' para ver los resultados.\n", LOG_FILE)
}

func main() {
	// 1. se abre el archivo para guardar los resultados
	logFile := createFileConnection(LOG_FILE, time.Now().Format(time.RFC3339))
	logFile.WriteString("Iteración | Latencia (ms) | Resultado\n")
	defer logFile.Close()

	// 2. Se conecta al servidor
	conn := createConnection(CONN_HOST, CONN_PORT, CONN_TYPE)
	defer conn.Close()

	// 3. Bucle de prueba de latencia
	totalLatency := runLatencyTest(conn, NUM_REQUESTS, logFile)

	// 4. Cálculo de latencia promedio
	createFinalReport(logFile, totalLatency)
}
