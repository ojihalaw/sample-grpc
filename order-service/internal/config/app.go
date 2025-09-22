package config

import (
	"fmt"
	"log"
	"net"

	"github.com/cloudinary/cloudinary-go/v2"
	controller "github.com/ojihalaw/shopping-cart-go-grpc/order-service/internal/delivery/grpc"
	"github.com/ojihalaw/shopping-cart-go-grpc/order-service/internal/repository"
	"github.com/ojihalaw/shopping-cart-go-grpc/order-service/internal/usecase"
	orderPb "github.com/ojihalaw/shopping-cart-go-grpc/shared/pb/order"
	utilsShared "github.com/ojihalaw/shopping-cart-go-grpc/shared/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB         *gorm.DB
	Log        *logrus.Logger
	Validator  *utilsShared.Validator
	Config     *viper.Viper
	Cloudinary *cloudinary.Cloudinary
}

func BootstrapGRPC(config *BootstrapConfig) {
	orderRepository := repository.NewOrderRepository(config.Log)
	orderUsecase := usecase.NewOrderUseCase(config.DB, config.Log, config.Validator, config.Cloudinary, orderRepository)
	orderController := controller.NewOrderController(orderUsecase, config.Log)

	port := config.Config.GetInt("APP_PORT")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	orderPb.RegisterOrderServiceServer(
		grpcServer,
		orderController,
	)
	message := fmt.Sprintf("✅ Order service running at : %d", port)
	config.Log.Println(message)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
