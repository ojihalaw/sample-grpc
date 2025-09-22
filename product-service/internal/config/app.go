package config

import (
	"fmt"
	"log"
	"net"

	"github.com/cloudinary/cloudinary-go/v2"
	controller "github.com/ojihalaw/shopping-cart-go-grpc/product-service/internal/delivery/grpc"
	"github.com/ojihalaw/shopping-cart-go-grpc/product-service/internal/repository"
	"github.com/ojihalaw/shopping-cart-go-grpc/product-service/internal/usecase"
	productpb "github.com/ojihalaw/shopping-cart-go-grpc/shared/pb/product"
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
	productRepository := repository.NewProductRepository(config.Log)
	productUsecase := usecase.NewProductUseCase(config.DB, config.Log, config.Validator, config.Cloudinary, productRepository)
	productController := controller.NewProductController(productUsecase, config.Log)

	port := config.Config.GetInt("APP_PORT")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	productpb.RegisterProductServiceServer(
		grpcServer,
		productController,
	)

	config.Log.Println("✅ Product service running at :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
