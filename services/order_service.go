package services

import (
	"crypto/sha512"
	"ecommerce-backend/dto"
	"ecommerce-backend/repositories"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type CheckoutResponse struct {
	OrderID    string `json:"orderId"`
	PaymentURL string `json:"paymentUrl"`
}

type OrderService interface {
	CheckoutCart(userID string) (CheckoutResponse, error)
	ProcessMidtransWebhook(payload map[string]interface{}) error

	GetAllOrdersAdmin() ([]dto.AdminOrderResponse, error)
	InputTrackingNumber(orderID string, req dto.TrackingRequest) error
}

type orderService struct {
	repo repositories.OrderRepository
}

func NewOrderService(repo repositories.OrderRepository) OrderService {
	return &orderService{repo}
}

func (s *orderService) CheckoutCart(userID string) (CheckoutResponse, error) {
	orderID, grossAmount, err := s.repo.Checkout(userID)
	if err != nil {
		return CheckoutResponse{}, err
	}

	var snapClient = snap.Client{}
	snapClient.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(grossAmount),
		},
	}

	snapResp, midtransErr := snapClient.CreateTransaction(req)

	if midtransErr != nil {
		fmt.Println("Error dari Midtrans:", midtransErr.Message)
		return CheckoutResponse{OrderID: orderID, PaymentURL: ""}, nil
	}

	s.repo.UpdatePaymentURL(orderID, snapResp.RedirectURL)

	return CheckoutResponse{
		OrderID:    orderID,
		PaymentURL: snapResp.RedirectURL,
	}, nil
}

func (s *orderService) ProcessMidtransWebhook(payload map[string]interface{}) error {
	orderID, _ := payload["order_id"].(string)
	statusCode, _ := payload["status_code"].(string)
	grossAmount, _ := payload["gross_amount"].(string)
	signatureKey, _ := payload["signature_key"].(string)
	transactionStatus, _ := payload["transaction_status"].(string)
	transactionID, _ := payload["transaction_id"].(string)
	paymentType, _ := payload["payment_type"].(string)

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	hashInput := orderID + statusCode + grossAmount + serverKey

	hasher := sha512.New()
	hasher.Write([]byte(hashInput))
	expectedSignature := hex.EncodeToString(hasher.Sum(nil))

	if expectedSignature != signatureKey {
		return fmt.Errorf("invalid signature key")
	}

	var newStatus string
	switch transactionStatus {
	case "capture", "settlement":
		newStatus = "Paid"
	case "deny", "cancel", "expire":
		newStatus = "Canceled"
	case "pending":
		newStatus = "Unpaid"
	default:
		newStatus = "Unknown"
	}

	return s.repo.UpdateOrderStatus(orderID, newStatus, transactionID, paymentType)
}

func (s *orderService) GetAllOrdersAdmin() ([]dto.AdminOrderResponse, error) {
	return s.repo.GetAllOrders()
}

func (s *orderService) InputTrackingNumber(orderID string, req dto.TrackingRequest) error {
	return s.repo.UpdateTrackingNumber(orderID, req.TrackingNumber)
}
