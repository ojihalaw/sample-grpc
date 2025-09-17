package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/ojihalaw/sample-grpc/product-service/internal/model"
	"github.com/ojihalaw/sample-grpc/product-service/internal/usecase"
	productpb "github.com/ojihalaw/sample-grpc/product-service/proto/product"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sirupsen/logrus"
)

type ProductController struct {
	Log     *logrus.Logger
	UseCase *usecase.ProductUseCase
	productpb.UnimplementedProductServiceServer
}

func NewProductController(useCase *usecase.ProductUseCase, logger *logrus.Logger) *ProductController {
	return &ProductController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *ProductController) GetProducts(ctx context.Context, req *productpb.GetProductsRequest) (*productpb.GetProductsResponse, error) {
	c.Log.Info("GetProducts called")

	// sementara dummy data
	products := []struct {
		ID   string
		Name string
	}{
		{"1", "Product A"},
		{"2", "Product B"},
	}

	var pbProducts []*productpb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &productpb.Product{
			Id:   p.ID,
			Name: p.Name,
		})
	}

	return &productpb.GetProductsResponse{
		Products: pbProducts,
	}, nil
}

func (c *ProductController) Create(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.ProductResponse, error) {
	c.Log.Info("GetProducts called")
	// Validasi field
	if req.Name == "" || req.Slug == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name and slug are required")
	}
	if req.Price < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "price must be >= 0")
	}
	if req.Stock < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "stock must be >= 0")
	}
	if req.CategoryId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "category_id is required")
	}

	// Convert categoryID ke UUID
	categoryID, err := uuid.Parse(req.CategoryId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid category_id")
	}

	// Panggil UseCase
	modelReq := &model.CreateProductRequest{
		Name:        req.Name,
		Slug:        req.Slug,
		SKU:         req.Sku,
		Variant:     req.Variant,
		Price:       int(req.Price),
		Stock:       int(req.Stock),
		Description: req.Description,
		CategoryID:  categoryID,
	}

	err = c.UseCase.Create(ctx, modelReq) // file upload via gRPC streaming atau terpisah
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	// Return response gRPC
	return &productpb.ProductResponse{
		Name: modelReq.Name,
		// Message: "product created successfully",
	}, nil
}
