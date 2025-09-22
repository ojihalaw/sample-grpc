package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/ojihalaw/sample-grpc/order-service/internal/entity"
	"github.com/ojihalaw/sample-grpc/order-service/internal/model"
	"github.com/ojihalaw/sample-grpc/order-service/internal/repository"
	utilsShared "github.com/ojihalaw/sample-grpc/shared/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderUseCase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validator       *utilsShared.Validator
	Cloudinary      *cloudinary.Cloudinary
	OrderRepository *repository.OrderRepository
}

func NewOrderUseCase(db *gorm.DB, logger *logrus.Logger, validator *utilsShared.Validator, cloudinary *cloudinary.Cloudinary, orderRepository *repository.OrderRepository) *OrderUseCase {
	return &OrderUseCase{
		DB:              db,
		Log:             logger,
		Validator:       validator,
		Cloudinary:      cloudinary,
		OrderRepository: orderRepository,
	}
}

func (o *OrderUseCase) Create(ctx context.Context, request *model.CreateOrderRequest) (*model.OrderResponse, error) {
	tx := o.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := o.Validator.Validate.Struct(request)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {

			var messages []string
			for _, e := range validationErrors {
				messages = append(messages, e.Translate(o.Validator.Translator))
			}
			return nil, fmt.Errorf("%w: %s", utilsShared.ErrValidation, strings.Join(messages, ", "))
		}
		return nil, fmt.Errorf("%w: %s", utilsShared.ErrValidation, err.Error())
	}

	todayCount, _ := o.OrderRepository.GetTodayOrderCount(ctx, tx)
	invoiceNumber := utilsShared.GenerateInvoice(todayCount + 1)

	var totalAmount int64
	var orderItems []entity.OrderItem
	for _, item := range request.Items {
		subtotal := item.Price * int64(item.Quantity)
		totalAmount += subtotal

		orderItems = append(orderItems, entity.OrderItem{
			ProductID: utilsShared.MustParseUUID(item.ProductID),
			Qty:       item.Quantity,
			Price:     item.Price,
			Subtotal:  subtotal,
		})
	}

	order := &entity.Order{
		UserID:        utilsShared.MustParseUUID(request.CustomerID),
		InvoiceNumber: invoiceNumber, // ex: INV-20250909-0003
		Status:        "Pending",     // pending / settlement / cancel
		Amount:        totalAmount,
		PaymentMethod: request.PaymentMethod, // ex: "e-wallet"
		PaymentType:   "resp.PaymentType",    // ex: "gopay"
		TransactionID: "resp.TransactionID",
		RedirectURL:   "deeplink", // kalau pakai coreapi
		OrderItems:    orderItems,
		Notes:         request.Notes,           // simpan catatan order
		ShippingAddr:  request.ShippingAddress, // simpan alamat pengiriman
	}
	if err := o.OrderRepository.Create(tx, order); err != nil {
		o.Log.Warnf("Failed create order to database : %+v", err)
		return nil, fmt.Errorf("%w: %s", utilsShared.ErrInternal, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		o.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fmt.Errorf("%w: %s", utilsShared.ErrInternal, err.Error())
	}

	return &model.OrderResponse{
		ID: order.ID.String(),
	}, nil
}
