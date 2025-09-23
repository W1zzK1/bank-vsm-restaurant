package main

import (
	"bank-vsm-restaurant/app/payment"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bank API is running!")
	})

    http.HandleFunc("/api/v1/payment",  payment.PaymentHandlerWithHTML)

    http.HandleFunc("/payment/success", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]int32{"payment_status": 200} // Success
		if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}})

	http.HandleFunc("/payment/fail", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]int32{"payment_status": 400} // Decline or fail
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}})
    
    log.Println("Server is starting on port :8181...")
	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
