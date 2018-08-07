package handler

import (
	"context"
	"server/common"
	"server/product-srv/db"
	product_proto "server/product-srv/proto/product"

	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

type ProductService struct{}

func (p *ProductService) AllProducts(ctx context.Context, req *product_proto.AllProductsRequest, rsp *product_proto.AllProductsResponse) error {
	log.Info("Received Product.All request")
	products, err := db.AllProducts(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(products) == 0 || err != nil {
		return common.NotFound(common.ProductSrv, p.AllProducts, err, "product not found")
	}
	rsp.Data = &product_proto.ProductArrData{products}
	return nil
}

func (p *ProductService) CreateProduct(ctx context.Context, req *product_proto.CreateProductRequest, rsp *product_proto.CreateProductResponse) error {
	log.Info("Received Product.Create request")
	if len(req.Product.Name) == 0 {
		return common.InternalServerError(common.ProductSrv, p.CreateProduct, nil, "product name empty")
	}
	if len(req.Product.Id) == 0 {
		req.Product.Id = uuid.NewUUID().String()
	}

	err := db.CreateProduct(ctx, req.Product)
	if err != nil {
		return common.InternalServerError(common.ProductSrv, p.CreateProduct, err, "create error")
	}
	rsp.Data = &product_proto.ProductData{req.Product}
	return nil
}

func (p *ProductService) ReadProduct(ctx context.Context, req *product_proto.ReadProductRequest, rsp *product_proto.ReadProductResponse) error {
	log.Info("Received Product.ReadProduct request")
	product, err := db.ReadProduct(ctx, req.Id, req.OrgId, req.TeamId)
	if product == nil || err != nil {
		return common.NotFound(common.ProductSrv, p.ReadProduct, err, "product not found")
	}
	rsp.Data = &product_proto.ProductData{product}
	return nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, req *product_proto.DeleteProductRequest, rsp *product_proto.DeleteProductResponse) error {
	log.Info("Received Product.DeleteProduct request")
	if err := db.DeleteProduct(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.ProductSrv, p.DeleteProduct, nil, "delete error")
	}
	return nil
}

func (p *ProductService) AutocompleteProduct(ctx context.Context, req *product_proto.AutocompleteProductRequest, rsp *product_proto.AutocompleteProductResponse) error {
	log.Info("Received Product.AutocompleteProduct request")
	products, err := db.AutocompleteProduct(ctx, req.Title)
	if len(products) == 0 || err != nil {
		return common.NotFound(common.ProductSrv, p.AutocompleteProduct, err, "product not found")
	}
	rsp.Data = &product_proto.ProductArrData{products}
	return nil
}

func (p *ProductService) AllServices(ctx context.Context, req *product_proto.AllServicesRequest, rsp *product_proto.AllServicesResponse) error {
	log.Info("Received Product.All request")
	services, err := db.AllServices(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(services) == 0 || err != nil {
		return common.NotFound(common.ProductSrv, p.AllServices, err, "service not found")
	}
	rsp.Data = &product_proto.ServiceArrData{services}
	return nil
}

func (p *ProductService) CreateService(ctx context.Context, req *product_proto.CreateServiceRequest, rsp *product_proto.CreateServiceResponse) error {
	log.Info("Received Product.Create request")
	if len(req.Service.Name) == 0 {
		return common.InternalServerError(common.ProductSrv, p.CreateService, nil, "service name empty")
	}
	if len(req.Service.Id) == 0 {
		req.Service.Id = uuid.NewUUID().String()
	}

	err := db.CreateService(ctx, req.Service)
	if err != nil {
		return common.InternalServerError(common.ProductSrv, p.CreateService, err, "create error")
	}
	rsp.Data = &product_proto.ServiceData{req.Service}
	return nil
}

func (p *ProductService) ReadService(ctx context.Context, req *product_proto.ReadServiceRequest, rsp *product_proto.ReadServiceResponse) error {
	log.Info("Received Product.ReadService request")
	service, err := db.ReadService(ctx, req.Id, req.OrgId, req.TeamId)
	if service == nil || err != nil {
		return common.NotFound(common.ProductSrv, p.ReadService, err, "service not found")
	}
	rsp.Data = &product_proto.ServiceData{service}
	return nil
}

func (p *ProductService) DeleteService(ctx context.Context, req *product_proto.DeleteServiceRequest, rsp *product_proto.DeleteServiceResponse) error {
	log.Info("Received Product.DeleteService request")
	if err := db.DeleteService(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.ProductSrv, p.DeleteService, nil, "delete error")
	}
	return nil
}

func (p *ProductService) AutocompleteService(ctx context.Context, req *product_proto.AutocompleteServiceRequest, rsp *product_proto.AutocompleteServiceResponse) error {
	log.Info("Received Product.AutocompleteService request")
	services, err := db.AutocompleteService(ctx, req.Title)
	if len(services) == 0 || err != nil {
		return common.NotFound(common.ProductSrv, p.AutocompleteService, err, "service not found")
	}
	rsp.Data = &product_proto.ServiceArrData{services}
	return nil
}

func (p *ProductService) AllBatches(ctx context.Context, req *product_proto.AllBatchesRequest, rsp *product_proto.AllBatchesResponse) error {
	log.Info("Received Product.All request")
	batches, err := db.AllBatches(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(batches) == 0 || err != nil {
		return common.NotFound(common.ProductSrv, p.AllBatches, err, "batch not found")
	}
	rsp.Data = &product_proto.BatchArrData{batches}
	return nil
}

func (p *ProductService) CreateBatch(ctx context.Context, req *product_proto.CreateBatchRequest, rsp *product_proto.CreateBatchResponse) error {
	log.Info("Received Product.Create request")
	if len(req.Batch.Name) == 0 {
		return common.InternalServerError(common.ProductSrv, p.CreateBatch, nil, "batch name empty")
	}
	if len(req.Batch.Id) == 0 {
		req.Batch.Id = uuid.NewUUID().String()
	}

	err := db.CreateBatch(ctx, req.Batch)
	if err != nil {
		return common.InternalServerError(common.ProductSrv, p.CreateBatch, err, "batch create error")
	}
	rsp.Data = &product_proto.BatchData{req.Batch}
	return nil
}

func (p *ProductService) ReadBatch(ctx context.Context, req *product_proto.ReadBatchRequest, rsp *product_proto.ReadBatchResponse) error {
	log.Info("Received Product.ReadBatch request")
	batch, err := db.ReadBatch(ctx, req.Id, req.OrgId, req.TeamId)
	if batch == nil || err != nil {
		return common.NotFound(common.ProductSrv, p.ReadBatch, err, "batch not found")
	}
	rsp.Data = &product_proto.BatchData{batch}
	return nil
}

func (p *ProductService) DeleteBatch(ctx context.Context, req *product_proto.DeleteBatchRequest, rsp *product_proto.DeleteBatchResponse) error {
	log.Info("Received Product.DeleteBatch request")
	if err := db.DeleteBatch(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.ProductSrv, p.DeleteBatch, nil, "delete error")
	}
	return nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, req *product_proto.UpdateProductRequest, rsp *product_proto.UpdateProductResponse) error {
	log.Info("Received Product.UpdateProduct request")
	return nil
}

func (p *ProductService) UpdateService(ctx context.Context, req *product_proto.UpdateServiceRequest, rsp *product_proto.UpdateServiceResponse) error {
	log.Info("Received Product.UpdateService request")
	return nil
}

func (p *ProductService) UpdateBatch(ctx context.Context, req *product_proto.UpdateBatchRequest, rsp *product_proto.UpdateBatchResponse) error {
	log.Info("Received Product.UpdateBatch request")
	return nil
}
