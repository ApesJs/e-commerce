package controllers

import (
	"e-commerce/app/consts"
	"e-commerce/app/models"
	"encoding/json"
	"github.com/midtrans/midtrans-go/snap"
	"io"
	"net/http"
)

func (server *Server) Midtrans(w http.ResponseWriter, r *http.Request) {
	var paymentNotification models.MidtransNotification

	err := json.NewDecoder(r.Body).Decode(&paymentNotification)
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}(r.Body)

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
