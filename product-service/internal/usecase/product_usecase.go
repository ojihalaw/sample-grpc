package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ojihalaw/shopping-cart-go-grpc/product-service/internal/entity"
	"github.com/ojihalaw/shopping-cart-go-grpc/product-service/internal/model"
	"github.com/ojihalaw/shopping-cart-go-grpc/product-service/internal/model/converter"
	"github.com/ojihalaw/shopping-cart-go-grpc/product-service/internal/repository"
	utilsShared "github.com/ojihalaw/shopping-cart-go-grpc/shared/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validator         *utilsShared.Validator
	Cloudinary        *cloudinary.Cloudinary
	ProductRepository *repository.ProductRepository
}

func NewProductUseCase(db *gorm.DB, logger *logrus.Logger, validator *utilsShared.Validator, cloudinary *cloudinary.Cloudinary, productRepository *repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		DB:                db,
		Log:               logger,
		Validator:         validator,
		Cloudinary:        cloudinary,
		ProductRepository: productRepository,
	}
}

func (p *ProductUseCase) Create(ctx context.Context, request *model.CreateProductRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := p.Validator.Validate.Struct(request)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {

			var messages []string
			for _, e := range validationErrors {
				messages = append(messages, e.Translate(p.Validator.Translator))
			}
			return fmt.Errorf("%w: %s", utilsShared.ErrValidation, strings.Join(messages, ", "))
		}
		return fmt.Errorf("%w: %s", utilsShared.ErrValidation, err.Error())
	}

	// check duplicate
	exists, err := p.ProductRepository.ExistsByName(p.DB.WithContext(ctx), request.Name)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("%w: %s", utilsShared.ErrConflict, "product name already exist")
	}
	slug := utilsShared.GenerateSlug(request.Name)
	sku := utilsShared.GenerateSKU(request.Name)

	product := &entity.Product{
		Name:        request.Name,
		SKU:         sku,
		Variant:     request.Variant,
		Price:       request.Price,
		Description: request.Description,
		Stock:       request.Stock,
		Slug:        slug,
		ImageURL:    "imageURL",
		CategoryID:  request.CategoryID,
		Star:        request.Star,
	}

	if err := p.ProductRepository.Create(p.DB.WithContext(ctx), product); err != nil {
		p.Log.Warnf("Failed create product to database : %+v", err)
		return fmt.Errorf("%w: %s", utilsShared.ErrInternal, err.Error())
	}

	return nil
}

func (p *ProductUseCase) FindAll(ctx context.Context, pagination *utilsShared.PaginationRequest) ([]model.ProductResponse, *utilsShared.PaginationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var products []entity.Product

	total, err := p.ProductRepository.FindAll(p.DB.WithContext(ctx), &products, pagination)
	if err != nil {
		p.Log.Warnf("Failed find all category from database : %+v", err)
		return nil, nil, fmt.Errorf("%w: %s", utilsShared.ErrInternal, err.Error())
	}

	responses := make([]model.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = *converter.ProductToResponse(&product)
	}

	totalPage := int((total + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginationRes := &utilsShared.PaginationResponse{
		Page:      pagination.Page,
		Limit:     pagination.Limit,
		OrderBy:   pagination.OrderBy,
		SortBy:    pagination.SortBy,
		Search:    pagination.Search,
		TotalData: total,
		TotalPage: totalPage,
	}

	return responses, paginationRes, nil
}

func (p *ProductUseCase) FindByID(ctx context.Context, productID string) (*model.ProductResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var products *entity.Product

	product, err := p.ProductRepository.FindById(p.DB.WithContext(ctx), products, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Log.Infof("product not found, id=%s", productID)
			return nil, utilsShared.ErrNotFound
		}
		p.Log.Warnf("Failed find product from database : %+v", err)
		return nil, fmt.Errorf("%w: %s", utilsShared.ErrInternal, err.Error())
	}

	return converter.ProductToResponse(product), nil
}

func (p *ProductUseCase) Update(ctx context.Context, productID string, request *model.UpdateProductRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Cari product lama
	product := &entity.Product{}
	_, err := p.ProductRepository.FindById(p.DB.WithContext(ctx), product, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utilsShared.ErrNotFound
		}
		return fmt.Errorf("%w: %s", utilsShared.ErrInternal, err.Error())
	}

	// Update field hanya kalau ada input baru
	if request.Name != "" {
		product.Name = request.Name
		product.Slug = utilsShared.GenerateSlug(request.Name)
		product.SKU = utilsShared.GenerateSKU(request.Name)
	}
	if request.Variant != "" {
		product.Variant = request.Variant
	}
	if request.Description != "" {
		product.Description = request.Description
	}
	if request.Price != 0 {
		product.Price = request.Price
	}
	if request.Stock != 0 {
		product.Stock = request.Stock
	}
	if request.CategoryID != uuid.Nil {
		product.CategoryID = request.CategoryID
	}

	// check duplicate
	exists, err := p.ProductRepository.ExistsByName(p.DB.WithContext(ctx), request.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: %s", utilsShared.ErrConflict, "product name already exist")
	}

	err = p.ProductRepository.Update(p.DB.WithContext(ctx), product)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Log.Infof("product not found, id=%s", productID)
			return utilsShared.ErrNotFound
		}
		p.Log.Warnf("Failed find product from database : %+v", err)
		return fmt.Errorf("%w: %s", utilsShared.ErrInternal, err.Error())
	}

	return nil
}

func (p *ProductUseCase) Delete(ctx context.Context, productID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	product := &entity.Product{}
	_, err := p.ProductRepository.FindById(p.DB.WithContext(ctx), product, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Log.Infof("product not found, id=%s", productID)
			return utilsShared.ErrNotFound
		}
		p.Log.Warnf("Failed find product from database : %+v", err)
		return fmt.Errorf("%w: %s", utilsShared.ErrInternal, err.Error())
	}

	err = p.ProductRepository.Delete(p.DB.WithContext(ctx), product)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Log.Infof("product not found, id=%s", productID)
			return utilsShared.ErrNotFound
		}
		p.Log.Warnf("Failed find product from database : %+v", err)
		return fmt.Errorf("%w: %s", utilsShared.ErrInternal, err.Error())
	}

	return nil
}
