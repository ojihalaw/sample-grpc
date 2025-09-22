package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/ojihalaw/sample-grpc/order-service/internal/model"
	"github.com/ojihalaw/sample-grpc/order-service/internal/usecase"
	orderPb "github.com/ojihalaw/sample-grpc/shared/pb/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// productpb "github.com/ojihalaw/sample-grpc/shared/pb/product"
	"github.com/sirupsen/logrus"
)

type OrderController struct {
	Log     *logrus.Logger
	UseCase *usecase.OrderUseCase
	orderPb.UnimplementedOrderServiceServer
}

func NewOrderController(useCase *usecase.OrderUseCase, logger *logrus.Logger) *OrderController {
	return &OrderController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *OrderController) Create(ctx context.Context, req *orderPb.CreateOrderRequest) (*orderPb.SuccessResponse, error) {

	c.Log.Info("✅ Create order called")
	var items []model.OrderItemRequest
	for _, i := range req.Items {
		items = append(items, model.OrderItemRequest{
			ProductID: i.ProductId,
			Name:      i.Name,
			Price:     i.Price,
			Quantity:  int(i.Quantity),
		})
	}

	modelReq := &model.CreateOrderRequest{
		CustomerID:      req.CustomerId,
		CustomerName:    req.CustomerName,
		CustomerEmail:   req.CustomerEmail,
		CustomerPhone:   req.CustomerPhone,
		Notes:           req.Notes,
		PaymentMethod:   req.PaymentMethod,
		ShippingAddress: req.ShippingAddress,
		Items:           items,
	}

	orderCreated, err := c.UseCase.Create(ctx, modelReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	return &orderPb.SuccessResponse{
		Code:    200,
		Status:  true,
		Message: "Order created successfully",
		Result: &orderPb.OrderResponse{
			Id: orderCreated.ID,
		},
	}, nil
}

func (c *OrderController) StreamOrderStatus(req *orderPb.OrderStatusRequest, stream orderPb.OrderService_StreamOrderStatusServer) error {
	c.Log.Info("✅ Stream Order Status called")
	orderID := req.GetOrderId()

	statuses := []string{"PLACED", "PROCESSED", "DELIVERED", "RECEIVED"}
	for _, st := range statuses {
		time.Sleep(5 * time.Second) // simulasi proses order
		res := &orderPb.OrderStatusResponse{
			OrderId:   orderID,
			Status:    st,
			Message:   fmt.Sprintf("order %s", st),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}

	return nil
}
