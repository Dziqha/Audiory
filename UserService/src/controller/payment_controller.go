package controller

import (
	"UserService/configs"
	"UserService/src/models"
	"UserService/src/res"
	"fmt"
	"log"
	"os"

	midtrans "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// MidtransPayment handles payment for subscriptions (daily, monthly, yearly)
func MidtransPayment(userID int, id string,price int,  subscriptionType string) (*res.PaymentResponse, error) {
	var customer models.User
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		log.Fatal("MIDTRANS_SERVER_KEY tidak disetel dalam environment variables")
	}

	midtrans.ServerKey = serverKey
	midtrans.Environment = midtrans.Sandbox

	// Ambil data pelanggan dari database
	if err := configs.Database().First(&customer, "id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil data pelanggan: %v", err)
	}
	var total int
	// Mengatur harga berdasarkan jenis langganan
	switch subscriptionType {
	case "day":
		price = price * 1
		total = price * 1 // Hitung total untuk beberapa hari
	case "month":
		price = price * 30 // Harga untuk sebulan (30 hari)
		total = price * 1 // Hitung total untuk beberapa bulan
	case "year":
		price = price * 365 // Harga untuk setahun (365 hari)
		total = price * 1 // Hitung total untuk beberapa tahun
	default:
		return nil, fmt.Errorf("jenis langganan tidak valid: %s", subscriptionType)
	}

	// Buat request transaksi ke Midtrans Snap
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  id,
			GrossAmt: int64(total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: customer.Username,
			Email: customer.Email,
			BillAddr: &midtrans.CustomerAddress{
				FName:    customer.Username,
				Phone:    customer.Email,
			},
		},
	}

	// Buat transaksi Snap
	snapResp, err := snap.CreateTransaction(req)
	if err != nil {
		log.Fatalf("Gagal membuat transaksi Snap: %v", err)
	}

	res := res.PaymentResponse{
		Token:       snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
	}

	fmt.Println("Result:", res)

	return &res, nil
}
