package controllers

import (
	"crypto/sha512"
	"e-commerce/app/consts"
	"e-commerce/app/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"os"
)

func (server *Server) Midtrans(w http.ResponseWriter, r *http.Request) {
	var paymentNotification models.MidtransNotification

	err := json.NewDecoder(r.Body).Decode(&paymentNotification)
	if err != nil {
		//todo: buatkan menjadi function agar tidak DRY
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		res := Result{
			Code:    http.StatusBadRequest,
			Data:    nil,
			Message: err.Error(),
		}
		response, _ := json.Marshal(res)

		_, _ = w.Write(response)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}(r.Body)

	fmt.Println("PaymentNotification nya nihhh", paymentNotification)

	err = validateSignatureKey(&paymentNotification)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)

		res := Result{
			Code:    http.StatusForbidden,
			Data:    nil,
			Message: err.Error(),
		}
		response, _ := json.Marshal(res)

		_, _ = w.Write(response)
		return
	}

	orderModel := models.Order{}
	order, err := orderModel.FindByID(server.DB, paymentNotification.OrderID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)

		res := Result{
			Code:    http.StatusForbidden,
			Data:    nil,
			Message: err.Error(),
		}
		response, _ := json.Marshal(res)

		_, _ = w.Write(response)
		return
	}

	if order.IsPaid() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)

		res := Result{
			Code:    http.StatusForbidden,
			Data:    nil,
			Message: "Payment sudah dilakukan",
		}
		response, _ := json.Marshal(res)

		_, _ = w.Write(response)
		return
	}

	paymentModel := models.Payment{}
	amount, _ := decimal.NewFromString(paymentNotification.GrossAmount)
	jsonPayload, _ := json.Marshal(paymentNotification)
	payload := (*json.RawMessage)(&jsonPayload)
	_, err = paymentModel.CreatePayment(server.DB, &models.Payment{
		OrderID:           order.ID,
		Amount:            amount,
		TransactionID:     paymentNotification.TransactionID,
		TransactionStatus: paymentNotification.TransactionStatus,
		Payload:           payload,
		PaymentType:       paymentNotification.PaymentType,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		res := Result{
			Code:    http.StatusBadRequest,
			Data:    nil,
			Message: err.Error(),
		}
		response, _ := json.Marshal(res)

		_, _ = w.Write(response)
		return
	}

	if IsPaymentSuccess(&paymentNotification) {
		err = order.MarkAsPaid(server.DB)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			res := Result{
				Code:    http.StatusBadRequest,
				Data:    nil,
				Message: "Payment tidak bisa diproses, payment gagal",
			}
			response, _ := json.Marshal(res)

			_, _ = w.Write(response)
			return
		}
		//todo: kirim receipt
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := Result{
		Code:    http.StatusOK,
		Data:    nil,
		Message: "Payment Saved",
	}
	response, _ := json.Marshal(res)

	_, _ = w.Write(response)
}

func IsPaymentSuccess(payload *models.MidtransNotification) bool {
	paymentStatus := false
	if payload.PaymentType == string(snap.PaymentTypeCreditCard) {
		paymentStatus = (payload.TransactionStatus == consts.PaymentStatusCapture) && (payload.FraudStatus == consts.FraudStatusAccept)
	} else {
		paymentStatus = (payload.TransactionStatus == consts.PaymentStatusSettlement) && (payload.FraudStatus == consts.FraudStatusAccept)
	}

	return paymentStatus
}

// validateSignatureKey will validate the signature key  in the midtrans payload
func validateSignatureKey(payload *models.MidtransNotification) error {
	//jika di production maka APP_ENV ubah menjadi PRODUCTION
	environment := os.Getenv("APP_ENV")
	if environment == "DEVELOPMENT" {
		return nil
	}

	signaturePayload := payload.OrderID + payload.StatusCode + payload.GrossAmount + os.Getenv("API_MIDTRANS_SERVER_KEY")
	sha512Value := sha512.New()
	sha512Value.Write([]byte(signaturePayload))

	signatureKey := fmt.Sprintf("%x", sha512Value.Sum(nil))

	if signatureKey != payload.SignatureKey {
		return errors.New("invalid signature key")
	}

	return nil
}
