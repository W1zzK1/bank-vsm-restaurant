package main

import (
    "bank-vsm-restaurant/app/payment"
    "net/http"
    "log"
    "fmt"
)

func main() {

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bank API is running!")
	})
      
    http.HandleFunc("/api/v1/payments", payment.PaymentHandler)
    
    log.Println("Server is starting on port :8181...")
	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
