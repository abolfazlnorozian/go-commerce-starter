package adminroute

import (
	"shop/admin/service"
	"shop/user/services"

	"github.com/gin-gonic/gin"
)

func AdminProductRoute(r *gin.RouterGroup) {
	p := r.Group("/")

	pro := r.Group("/")
	count := r.Group("/")

	count.GET("/products/count", services.GetAllCountProducts)

	p.GET("/products", service.GetAllProductsByAdmin)

	pro.POST("/products", service.PostProductByAdmin)
	pro.DELETE("/products/:id", service.DeleteProductByAdmin)

	pro.GET("/products/byid/:id", service.CommentGetProductByIDByAdmin)
	//*************************************************
	p.POST("/products/:id/dimensions", service.PostDimensionByAdmin)

	p.POST("/products/:id/variations/", service.PostVariationToProductByAdmin)

	//*******//**************************
	p.GET("/products/:id", func(c *gin.Context) {
		values := c.Query("values")

		// Check if 'values' is 'false'
		if values == "false" {
			// Call the handler for GetProductByIDByValueFalse
			service.GetProductByIDByVAlueFalse(c)
			return
		}

		// Otherwise, call the handler for GetProductByIdByAdmin
		service.GetProductByIdByAdmin(c)
	})

	p.PUT("/products/:id", func(c *gin.Context) {
		// Check if there's a variation ID in the request URL
		if variationId := c.Param("variationId"); variationId != "" {

		} else {
			// If no variation ID is present, update the product
			service.UpdateProductByAdmin(c)
		}
	})

	p.PUT("/products/:id/variations/:variationId", func(c *gin.Context) {

		service.UpdateVariationByAdmin(c)
	})

	pro.DELETE("/products/:id/variations/:variatId", service.DeleteVariationFromProductByAdmin)

	p.POST("/products/:id/dimensions/:dimensionKey/values", service.PostValuesByAdminToDimension)

	p.DELETE("/products/:id/dimensions/:keyID", service.DeleteDimensionFromProductByAdmin)

	p.GET("/mix-products", service.GetMixProductsByAdmin)

	pro.POST("/mix-products", service.PostMixByAdmin)

	pro.DELETE("/mix-products/:id", service.DeleteMixBYAdmin)

}

func AdminRoutes(r *gin.RouterGroup) {
	u := r.Group("/admins")

	u.POST("/login", service.LoginAdmin)

}

func AdminCategoryRouter(r *gin.RouterGroup) {

	ca := r.Group("/")

	categ := r.Group("/")

	ca.GET("/categories", service.FindCategoriesByAdmin)

	categ.POST("/categories", service.PostCategoriesByAdmin)
	categ.DELETE("/categories/:id", service.DeleteCategoryByAdmin)

	categ.PUT("/categories/:id", service.UpdateCategoryByAdmin)

}
func AdminBrandRoute(r *gin.RouterGroup) {
	b := r.Group("/")
	b2 := r.Group("/")

	b.GET("/brands/", service.GetBrandsByAdmin)

	b2.POST("/brands", service.PostBrandsByAdmin)

	b2.DELETE("/brands/:id", service.DeleteBrandsByAdmin)

	b.PUT("/brands/:id", service.UpdateBrandByAdmin)

}

func Comments(r *gin.RouterGroup) {
	co := r.Group("/")

	co.GET("/comments", service.GetCommentByAdmin)

	co.PUT("/comments/:id", service.UpdateCommentByAdmin)

}

func UserRouteByAdmin(r *gin.RouterGroup) {
	co := r.Group("/")

	us := r.Group("/users")
	co.GET("/users", service.GetUsersByAdmin)

	us.GET("/count", service.GetAllCountUsers)

}

func PropertiesByAdmin(r *gin.RouterGroup) {

	//********************************************************
	pe := r.Group("/products")
	//******************************************************
	p := r.Group("/products")

	p.GET("/properties", func(c *gin.Context) {
		// Check if 'id' query parameter exists
		id := c.Query("id")
		if id != "" {
			service.GetPropertiesByIdByAdmin(c)

			return
		}

		// If 'id' parameter is not provided, get all properties
		service.GetPropertiesByAdmin(c)
	})

	// p.GET("/properties", service.GetPropertiesById)
	pe.POST("/properties", service.PostPropertiesByAdmin)

	pe.DELETE("/properties/:id", service.DeletePropertiesByAdmin)

	pe.PUT("/properties/:id", service.UpdatePropertiesByAdmin)

}
func ImageRouteByAdmin(r *gin.RouterGroup) {
	img := r.Group("/")

	img.GET("/uploads", service.FindAllImagesByAdmin)

	img.POST("/uploads", service.UploadImageByAdmin)
	img.DELETE("/uploads/:id", service.DeleteImageByAdmin)
	// down.GET("/uploads/:filename", service.FindOneImage)
	// down := r.Group("/")
	// down.Static("/uploads", "./public/images")
}
func DownloaderByAdmin(r *gin.RouterGroup) {
	down := r.Group("/")
	down.Static("/uploads", "./public/images")
	//down.GET("/uploads/:filename", upload.FindOneImage)

}
func PageRouteByAdmin(r *gin.RouterGroup) {
	pg := r.Group("/")

	pg.GET("/pages", service.GetPagesByAdmin)
	pg.POST("/pages", service.PostPagesByAdmin)

	pg.PUT("/pages/:id", service.UpdatePagesByAdmin)
}
func OrderRouteByAdmin(r *gin.RouterGroup) {
	or := r.Group("/")
	ord := r.Group("/")

	ord.GET("/users/orders", service.GetOrdersByAdmin)

	ord.GET("/users/orders/:id", service.GetOrderByIdByAdmin)
	ord.GET("/users/:id/information", service.GetInformationUserByAdmin)
	ord.POST("/users/orders/:id/status", service.UpdateSatusOrderByAdmin)
	or.GET("/users/orders/count", service.GetAllOrders)
}
func CouponRoutesByAdmin(r *gin.RouterGroup) {
	coupon := r.Group("/")

	coupon.GET("/coupons", service.GetCouponsByAdmin)

	coupon.POST("/coupons", service.PostCopounByAdmin)
	coupon.DELETE("/coupons/:id", service.DeleteCouponByAdmin)
	coupon.PUT("/coupons/:id", service.UpdateCouponByAdmin)
}
