package converter

import (
	"github.com/ojihalaw/sample-grpc/order-service/internal/entity"
	"github.com/ojihalaw/sample-grpc/order-service/internal/model"
)

func OrderToResponse(order *entity.Order) *model.OrderResponse {
	return &model.OrderResponse{
		ID: order.ID.String(),
	}
}
