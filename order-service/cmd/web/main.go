package main

import (
	"github.com/ojihalaw/shopping-cart-go-grpc/order-service/internal/config"
	configShared "github.com/ojihalaw/shopping-cart-go-grpc/shared/config"
	utilsShared "github.com/ojihalaw/shopping-cart-go-grpc/shared/utils"
)

func main() {
	viperConfig := config.NewViper()
	log := configShared.NewLogger(viperConfig)
	db := configShared.NewDatabase(viperConfig, log)
	validator := utilsShared.NewValidator(viperConfig)
	cloudinary := configShared.NewCloudinary(viperConfig)

	bootstrapConfig := &config.BootstrapConfig{
		DB:         db,
		Log:        log,
		Config:     viperConfig,
		Validator:  validator,
		Cloudinary: cloudinary,
	}

	config.BootstrapGRPC(bootstrapConfig)
}
