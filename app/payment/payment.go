package payment

import (
	"encoding/json"
	"log"
	"net/http"
	"html/template"
)


type Payment struct {
	OrderID  string   `json:"order_id"`
	Amount   float64  `json:"amount"`
}

type PageData struct {
	Payment Payment
}

func PaymentHandlerWithHTML(w http.ResponseWriter, r *http.Request) {
	// 1. Проверяем, что это POST-запрос
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method. Only POST is allowed.", http.StatusMethodNotAllowed)
		return
	}

	// 2. Декодируем JSON из тела запроса
	var req Payment
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Парсим наш HTML-шаблон
	tmpl, err := template.ParseFiles("web/payment.html")
	if err != nil {
		http.Error(w, "Could not parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Создаем данные для передачи в шаблон
	data := PageData{
		Payment: req,
	}

	// 5. "Исполняем" шаблон: вставляем данные `data` в шаблон `tmpl`
	// и отправляем результат пользователю `w`.
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Could not execute template: "+err.Error(), http.StatusInternalServerError)
	}
}

// PaykeeperWebhookHandler принимает серверные уведомления от Paykeeper.
func PaykeeperWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Вебхуки всегда приходят методом POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 2. Paykeeper отправляет данные в формате x-www-form-urlencoded, а не JSON
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Cannot parse form", http.StatusBadRequest)
		return
	}

	// 3. Получаем ключевые поля из формы
	orderID := r.FormValue("orderid")
	status := r.FormValue("status") // "paid" или "failed"

	log.Printf("Webhook received for order %s with status '%s'", orderID, status)

	// 5. Обновляем статус заказа в вашей базе данных
	if status == "paid" {
		log.Printf("SUCCESS: Order %s is paid.", orderID)
		// db.UpdateOrderStatus(orderID, "paid")
	} else {
		log.Printf("FAILURE: Order %s is not paid (status: %s).", orderID, status)
		// db.UpdateOrderStatus(orderID, "failed")
	}

	// 6. Отвечаем Paykeeper, что мы получили вебхук.
	// Если не ответить 200 OK, он будет пытаться отправить его снова.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}