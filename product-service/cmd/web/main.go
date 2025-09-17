package main

import (
	"github.com/ojihalaw/sample-grpc/common/utils"
	"github.com/ojihalaw/sample-grpc/product-service/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	app := config.NewFiber(viperConfig)
	validator := utils.NewValidator(viperConfig)
	cloudinary := config.NewCloudinary(viperConfig)

	bootstrapConfig := &config.BootstrapConfig{
		DB:         db,
		App:        app,
		Log:        log,
		Config:     viperConfig,
		Validator:  validator,
		Cloudinary: cloudinary,
	}

	config.BootstrapGRPC(bootstrapConfig)
}
