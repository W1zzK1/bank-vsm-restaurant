package payment

import (
	"errors"
	"encoding/json"
	"strings"
	"log"
	"net/http"
)

type CardData struct {
	CardNumber string `json:"card_number"`
	ExpMonth   int    `json:"exp_month"`
	ExpYear    int    `json:"exp_year"`
	CVV        string `json:"cvv"`
}

type Payment struct {
	OrderID        string            `json:"order_id"`
	Amount         float32           `json:"amount"`
	Currency       string            `json:"currency"`
	Description    string            `json:"description"`
	Card           CardData          `json:"card_data"`
	Metadata       map[string]string `json:"metadata"`
	IPAddress      string            `json:"ip_address"`
	IdempotencyKey string            `json:"idempotency_key"` /*
	это уникальный идентификатор, генерируемый клиентом и передаваемый в запросе на сервер.
	Он позволяет серверу распознавать повторные запросы и выполнять операцию только 
	один раз, предотвращая дублирование и нежелательные побочные эффекты, например, 
	двойное списание средств. 
	*/
}

func (p *Payment) Validate() error {
	if p.OrderID == "" {
		return errors.New("order_id is required")
	}
	if p.Amount <= 0 {
		return errors.New("amount must be a positive value")
	}
	if p.Currency == "" {
		return errors.New("currency is required")
	}
	if p.Card.CardNumber == "" {
		return errors.New("card_number is required")
	}
	// Здесь можно добавить больше проверок: формат CVV, срок действия карты и т.д.
	return nil
}

func (p *Payment) Authorize() bool {
	// Очищаем номер карты от пробелов для удобства сравнения.
	cardNumber := strings.ReplaceAll(p.Card.CardNumber, " ", "")

	switch cardNumber {
	case "4242424242424242":
		// Это стандартная тестовая карта для успешных платежей.
		return true
	case "4000000000000002":
		// Тестовая карта для получения отказа.
		return false
	default:
		// Для всех остальных карт в этом примере возвращаем отказ.
		return false
	}
}

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request received for: %s from %s", r.URL.Path, r.RemoteAddr)

	// Принимаем только POST запросы
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем JSON из тела запроса в нашу структуру Payment
	var p Payment
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// 1. Валидация данных
	if err := p.Validate(); err != nil {
		log.Printf("Validation failed for order %s: %v", p.OrderID, err)
		http.Error(w, err.Error(), http.StatusBadRequest) // 400 Bad Request
		return
	}

	// 2. Авторизация платежа
	if !p.Authorize() {
		log.Printf("Authorization declined for order %s (card: ...%s)", p.OrderID, p.Card.CardNumber[len(p.Card.CardNumber)-4:])
		http.Error(w, "Payment declined", http.StatusPaymentRequired) // 402 Payment Required
		return
	}

	// Если все успешно
	log.Printf("Payment successful for order %s", p.OrderID)
	w.WriteHeader(http.StatusOK) // 200 OK
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "order_id": p.OrderID})
}
