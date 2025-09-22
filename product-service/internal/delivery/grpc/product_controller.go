package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/ojihalaw/shopping-cart-go-grpc/product-service/internal/model"
	"github.com/ojihalaw/shopping-cart-go-grpc/product-service/internal/usecase"
	productpb "github.com/ojihalaw/shopping-cart-go-grpc/shared/pb/product"
	utilsShared "github.com/ojihalaw/shopping-cart-go-grpc/shared/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (c *ProductController) GetProducts(ctx context.Context, req *productpb.GetProductsRequest) (*productpb.SuccessResponseWithPagination, error) {
	c.Log.Info("✅ GetProducts called")

	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10 // default limit
	}

	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "created_at" // default order field
	}

	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "desc" // default sort
	}

	search := req.Search

	paginationReq := &utilsShared.PaginationRequest{
		Page:    int(page),
		Limit:   int(limit),
		OrderBy: orderBy,
		SortBy:  sortBy,
		Search:  search,
	}

	products, pagination, err := c.UseCase.FindAll(ctx, paginationReq)
	if err != nil {
		return &productpb.SuccessResponseWithPagination{
			Code:    500,
			Status:  false,
			Message: err.Error(),
		}, nil
	}

	var pbProducts []*productpb.ProductResponse
	for _, p := range products {
		pbProducts = append(pbProducts, &productpb.ProductResponse{
			Id:           p.ID,
			Name:         p.Name,
			Price:        int32(p.Price),
			Slug:         p.Slug,
			Sku:          p.SKU,
			Variant:      p.Variant,
			Stock:        int32(p.Stock),
			Description:  p.Description,
			Star:         p.Star,
			ImageUrl:     p.ImageURL,
			CategoryId:   p.CategoryID.String(),
			CategoryName: p.CategoryName,
			CreatedAt:    p.CreatedAt,
			UpdatedAt:    p.UpdatedAt,
		})
	}

	return &productpb.SuccessResponseWithPagination{
		Code:    200,
		Status:  true,
		Message: "Products fetched successfully",
		Result: &productpb.GetProductsResponse{
			Products: pbProducts,
		},
		Pagination: &productpb.Pagination{
			Page:  int32(pagination.Page),
			Limit: int32(pagination.Limit),
			Total: int32(pagination.TotalData),
		},
	}, nil
}

func (c *ProductController) FindById(ctx context.Context, req *productpb.GetProductByIDRequest) (*productpb.SuccessResponse, error) {
	c.Log.Info("✅ FindById product called")

	product, err := c.UseCase.FindByID(ctx, req.Id)
	if err != nil {
		return &productpb.SuccessResponse{
			Code:    500,
			Status:  false,
			Message: err.Error(),
			Result:  nil,
		}, nil
	}

	pbProduct := &productpb.ProductResponse{
		Id:           product.ID,
		Name:         product.Name,
		Price:        int32(product.Price),
		Slug:         product.Slug,
		Sku:          product.SKU,
		Variant:      product.Variant,
		Stock:        int32(product.Stock),
		Description:  product.Description,
		Star:         float64(product.Star),
		ImageUrl:     product.ImageURL,
		CategoryId:   product.CategoryID.String(),
		CategoryName: product.CategoryName,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}

	return &productpb.SuccessResponse{
		Code:    200,
		Status:  true,
		Message: "Detail product fetched successfully",
		Result:  pbProduct,
	}, nil
}

func (c *ProductController) Create(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.DefaultSuccessResponse, error) {

	c.Log.Info("✅ Create called")
	// Validasi field
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name are required")
	}
	if req.Price < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "price must be >= 0")
	}
	if req.Stock < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "stock must be >= 0")
	}
	if req.Variant == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Variant is required")
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
		Variant:     req.Variant,
		Price:       int(req.Price),
		Stock:       int(req.Stock),
		Description: req.Description,
		CategoryID:  categoryID,
		Star:        float64(req.Star),
	}

	err = c.UseCase.Create(ctx, modelReq) // file upload via gRPC streaming atau terpisah
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return &productpb.DefaultSuccessResponse{
		Code:    200,
		Status:  true,
		Message: "Product created successfully",
	}, nil
}
