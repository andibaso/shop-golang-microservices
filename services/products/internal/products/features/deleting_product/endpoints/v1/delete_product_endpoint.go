package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"net/http"

	"github.com/meysamhadeli/shop-golang-microservices/pkg/mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/deleting_product"
)

type deleteProductEndpoint struct {
	*config.ProductEndpointBase[config.InfrastructureConfiguration]
}

func NewDeleteProductEndpoint(productEndpointBase *config.ProductEndpointBase[config.InfrastructureConfiguration]) *deleteProductEndpoint {
	return &deleteProductEndpoint{productEndpointBase}
}

func (ep *deleteProductEndpoint) MapRoute() {
	ep.ProductsGroup.DELETE("/:id", ep.deleteProduct())
}

// DeleteProduct
// @Tags Products
// @Summary Delete product
// @Description Delete existing product
// @Accept json
// @Produce json
// @Success 204
// @Param id path string true "Product ID"
// @Router /api/v1/products/{id} [delete]
func (ep *deleteProductEndpoint) deleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		request := &deleting_product.DeleteProductRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Configuration.Log.Warn("Bind", err)
			return err
		}

		command := deleting_product.NewDeleteProduct(request.ProductID)

		if err := ep.Configuration.Validator.StructCtx(ctx, command); err != nil {
			ep.Configuration.Log.Warn("validate", err)
			return err
		}

		_, err := mediatr.Send[*mediatr.Unit](ctx, command)

		if err != nil {
			ep.Configuration.Log.Warn("DeleteProduct", err)
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}