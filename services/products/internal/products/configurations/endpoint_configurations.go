package configurations

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	creating_product "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/creating_product/endpoints/v1"
	deleting_product "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/deleting_product/endpoints/v1"
	gettting_product_by_id "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/getting_product_by_id/endpoints/v1"
	getting_products "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/getting_products/endpoints/v1"
	searching_products "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/searching_product/endpoints/v1"
	updating_product "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/updating_product/endpoints/v1"
)

func (c *productsModuleConfigurator) configEndpoints(ctx context.Context, group *echo.Group) {

	productEndpointBase := &config.ProductEndpointBase[config.InfrastructureConfiguration]{
		ProductsGroup: group,
		Configuration: *c.InfrastructureConfiguration,
	}
	// CreateNewProduct
	createProductEndpoint := creating_product.NewCreteProductEndpoint(productEndpointBase)
	createProductEndpoint.MapRoute()

	// UpdateProduct
	updateProductEndpoint := updating_product.NewUpdateProductEndpoint(productEndpointBase)
	updateProductEndpoint.MapRoute()

	// GetProducts
	getProductsEndpoint := getting_products.NewGetProductsEndpoint(productEndpointBase)
	getProductsEndpoint.MapRoute()

	// SearchProducts
	searchProducts := searching_products.NewSearchProductsEndpoint(productEndpointBase)
	searchProducts.MapRoute()

	// GetProductById
	getProductByIdEndpoint := gettting_product_by_id.NewGetProductByIdEndpoint(productEndpointBase)
	getProductByIdEndpoint.MapRoute()

	// DeleteProduct
	deleteProductEndpoint := deleting_product.NewDeleteProductEndpoint(productEndpointBase)
	deleteProductEndpoint.MapRoute()
}