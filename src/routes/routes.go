package routes

import (
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/controllers"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/middlewares"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/dreamskynl/godi"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(r *gin.Engine, sc godi.IGoDI) {
	// 健康检查端点（不需要认证）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	registerManagementRoutes(r.Group("/api/management"), sc)
	registerProductionRoutes(r.Group("/api/production"), sc)
}

func registerManagementRoutes(r *gin.RouterGroup, sc godi.IGoDI) {
	r.POST("/login", func(c *gin.Context) { controllers.NewManagementController(c, sc).Login() })

	// Authorized routes
	r.Use(middlewares.AuthorizeJWT(models.JwtServiceRoleAdmin))
	{
		r.POST("/supplier", func(c *gin.Context) { controllers.NewManagementController(c, sc).AddSupplier() })
		r.DELETE("/supplier", func(c *gin.Context) { controllers.NewManagementController(c, sc).DeleteSupplier() })
		r.GET("/supplier", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetSuppliers() })
		r.GET("/supplier/:id", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetSupplier() })
		r.PUT("/supplier", func(c *gin.Context) { controllers.NewManagementController(c, sc).UpdateSupplier() })

		r.POST("/product_model", func(c *gin.Context) { controllers.NewManagementController(c, sc).AddProductModel() })
		r.DELETE("/product_model", func(c *gin.Context) { controllers.NewManagementController(c, sc).DeleteProductModel() })
		r.GET("/product_model", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetProductModels() })
		r.GET("/product_model/:id", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetProductModel() })
		r.PUT("/product_model", func(c *gin.Context) { controllers.NewManagementController(c, sc).UpdateProductModel() })

		r.POST("/production_plan", func(c *gin.Context) { controllers.NewManagementController(c, sc).AddProductionPlan() })
		r.DELETE("/production_plan", func(c *gin.Context) { controllers.NewManagementController(c, sc).DeleteProductionPlan() })
		r.GET("/production_plan", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetProductionPlans() })
		r.GET("/production_plan/:id", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetProductionPlan() })
		r.PUT("/production_plan", func(c *gin.Context) { controllers.NewManagementController(c, sc).UpdateProductionPlan() })
		r.GET("/production_plan/date_range", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetProductionPlansByDateRange() })

		r.GET("/product_line", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetProductLines() })
		r.GET("/product_line/:id", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetProductLine() })
		r.POST("/product_line", func(c *gin.Context) { controllers.NewManagementController(c, sc).AddProductLine() })
		r.DELETE("/product_line", func(c *gin.Context) { controllers.NewManagementController(c, sc).DeleteProductLine() })

		r.GET("/pallet", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetPallets() })
		r.GET("/pallet/:id", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetPallet() })

		r.GET("/product", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetProducts() })
		r.GET("/product/:id", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetProduct() })

		r.POST("/api", func(c *gin.Context) { controllers.NewManagementController(c, sc).AddApi() })
		r.DELETE("/api", func(c *gin.Context) { controllers.NewManagementController(c, sc).DeleteApi() })
		r.GET("/api", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetApis() })
		r.GET("/api/:id", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetApi() })
		r.PUT("/api", func(c *gin.Context) { controllers.NewManagementController(c, sc).UpdateApi() })

		r.POST("/user", func(c *gin.Context) { controllers.NewManagementController(c, sc).AddUser() })
		r.DELETE("/user", func(c *gin.Context) { controllers.NewManagementController(c, sc).DeleteUser() })
		r.GET("/user", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetUsers() })
		r.GET("/user/:id", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetUser() })
		r.PUT("/user", func(c *gin.Context) { controllers.NewManagementController(c, sc).UpdateUser() })

		// 质量统计相关接口
		r.GET("/quality_stats", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetQualityStats() })

		// 数据报表相关接口
		r.GET("/report/defect", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetDefectReport() })
		r.GET("/report/inspection", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetInspectionReport() })
		r.GET("/report/cost", func(c *gin.Context) { controllers.NewManagementController(c, sc).GetCostReport() })
	}
}

func registerProductionRoutes(r *gin.RouterGroup, sc godi.IGoDI) {
	r.POST("/register", func(c *gin.Context) { controllers.NewProductionController(c, sc).RegisterProductLine() })
	r.POST("/authenticate", func(c *gin.Context) { controllers.NewProductionController(c, sc).AuthenticateProductLine() })

	// 需要产线认证的接口
	r.Use(middlewares.AuthorizeProductionLineJWT())
	{
		r.POST("/product_line", func(c *gin.Context) { controllers.NewProductionController(c, sc).AddProductLine() })
		r.DELETE("/product_line", func(c *gin.Context) { controllers.NewProductionController(c, sc).DeleteProductLine() })

		r.POST("/pallet", func(c *gin.Context) { controllers.NewProductionController(c, sc).AddPallet() })

		r.POST("/product", func(c *gin.Context) { controllers.NewProductionController(c, sc).AddProduct() })
	}
}
