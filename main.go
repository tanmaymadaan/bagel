package main

import (
	"fmt"
	"net/http"
)

func main() {
	servers := []string{
		"http://localhost:8081/",
		"http://localhost:8082/",
		"http://localhost:8083/",
		"http://localhost:8084/",
		"http://localhost:8085/",
	}

	loadBalancer := NewLoadBalancer(servers)

	http.Handle("/", loadBalancer)

	fmt.Println("Backend server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %s\n", err)
	}
}
