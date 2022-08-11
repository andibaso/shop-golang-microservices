package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"net/http"

	"github.com/meysamhadeli/shop-golang-microservices/pkg/mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/searching_product"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/searching_product/dtos"
)

type searchProductsEndpoint struct {
	*config.ProductEndpointBase[config.InfrastructureConfiguration]
}

func NewSearchProductsEndpoint(productEndpointBase *config.ProductEndpointBase[config.InfrastructureConfiguration]) *searchProductsEndpoint {
	return &searchProductsEndpoint{productEndpointBase}
}

func (ep *searchProductsEndpoint) MapRoute() {
	ep.ProductsGroup.GET("/search", ep.searchProducts())
}

// SearchProducts
// @Tags Products
// @Summary Search products
// @Description Search products
// @Accept json
// @Produce json
// @Param searchProductsRequestDto query dtos.SearchProductsRequestDto false "SearchProductsRequestDto"
// @Success 200 {object} dtos.SearchProductsResponseDto
// @Router /api/v1/products/search [get]
func (ep *searchProductsEndpoint) searchProducts() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		listQuery, err := utils.GetListQueryFromCtx(c)

		if err != nil {
			utils.LogResponseError(c, ep.Configuration.Log, err)
			return err
		}

		request := &dtos.SearchProductsRequestDto{ListQuery: listQuery}

		// https://echo.labstack.com/guide/binding/
		if err := c.Bind(request); err != nil {
			ep.Configuration.Log.Warn("Bind", err)
			return err
		}

		query := &searching_product.SearchProducts{SearchText: request.SearchText, ListQuery: request.ListQuery}

		if err := ep.Configuration.Validator.StructCtx(ctx, query); err != nil {
			ep.Configuration.Log.Errorf("(validate) err: {%v}", err)
			return err
		}

		queryResult, err := mediatr.Send[*dtos.SearchProductsResponseDto](ctx, query)

		if err != nil {
			ep.Configuration.Log.Warn("SearchProducts", err)
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}