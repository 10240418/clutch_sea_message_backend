package controllers

import (
	"time"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/services"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/utils"
	"github.com/dreamskynl/godi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ManagementController struct {
	ctx                   *gin.Context
	supplierService       services.ISupplierService
	productModelService   services.IProductModelService
	productionPlanService services.IProductionPlanService
	productLineService    services.IProductLineService
	palletService         services.IPalletService
	productService        services.IProductService
	apiService            services.IAPIService
	userService           services.IUserService
	jwtService            services.IJwtService
	qualityStatsService   services.IQualityStatsService
	dataReportService     services.IDataReportService
}

func NewManagementController(ctx *gin.Context, sc godi.IGoDI) IManagementController {
	return &ManagementController{
		ctx:                   ctx,
		supplierService:       sc.MustResolve(&services.SupplierService{}).(*services.SupplierService),
		productModelService:   sc.MustResolve(&services.ProductModelService{}).(*services.ProductModelService),
		productionPlanService: sc.MustResolve(&services.ProductionPlanService{}).(*services.ProductionPlanService),
		productLineService:    sc.MustResolve(&services.ProductLineService{}).(*services.ProductLineService),
		palletService:         sc.MustResolve(&services.PalletService{}).(*services.PalletService),
		productService:        sc.MustResolve(&services.ProductService{}).(*services.ProductService),
		apiService:            sc.MustResolve(&services.APIService{}).(*services.APIService),
		userService:           sc.MustResolve(&services.UserService{}).(*services.UserService),
		jwtService:            sc.MustResolve(&services.JwtService{}).(*services.JwtService),
		qualityStatsService:   sc.MustResolve(&services.QualityStatsService{}).(*services.QualityStatsService),
		dataReportService:     sc.MustResolve(&services.DataReportService{}).(*services.DataReportService),
	}
}

func (mc *ManagementController) AddSupplier() {
	var form models.Supplier
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.supplierService.CreateSupplier(&form); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(201, gin.H{"data": form, "message": "success"})
}

func (mc *ManagementController) DeleteSupplier() {
	var form IDsField
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.supplierService.DeleteSuppliers(form.IDs); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"message": "success"})
}

func (mc *ManagementController) GetSuppliers() {
	var queryParams struct{}
	var paginateParams models.PaginationQuery
	if err := mc.ctx.ShouldBindQuery(&queryParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.ctx.ShouldBindQuery(&paginateParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	queryParamsMap := utils.StructToMap(queryParams)
	paginateParamsMap := utils.StructToMap(paginateParams)
	suppliers, pageResult, err := mc.supplierService.GetSuppliers(queryParamsMap, paginateParamsMap)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": suppliers, "pagination": pageResult, "message": "success"})
}

func (mc *ManagementController) GetSupplier() {
	var uriParams IDField
	if err := mc.ctx.ShouldBindUri(&uriParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	supplier, err := mc.supplierService.GetSupplier(uriParams.ID)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": supplier, "message": "success"})
}

func (mc *ManagementController) UpdateSupplier() {
	var form models.Supplier
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	supplier, err := mc.supplierService.GetSupplier(int64(form.ID))
	if err != nil || supplier.ID == 0 {
		mc.ctx.JSON(404, gin.H{"error": "supplier not found"})
		return
	}

	supplierMap := utils.StructToMap(form)
	if err := mc.supplierService.UpdateSupplier(supplier, supplierMap); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": form, "message": "success"})
}

func (mc *ManagementController) AddProductModel() {
	var form models.ProductModel
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.productModelService.CreateProductModel(&form); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(201, gin.H{"data": form, "message": "success"})
}

func (mc *ManagementController) DeleteProductModel() {
	var form IDsField
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.productModelService.DeleteProductModels(form.IDs); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"message": "success"})
}

func (mc *ManagementController) GetProductModels() {
	var queryParams struct{}
	var paginateParams models.PaginationQuery
	if err := mc.ctx.ShouldBindQuery(&queryParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.ctx.ShouldBindQuery(&paginateParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	queryParamsMap := utils.StructToMap(queryParams)
	paginateParamsMap := utils.StructToMap(paginateParams)
	productModels, pageResult, err := mc.productModelService.GetProductModels(queryParamsMap, paginateParamsMap)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": productModels, "pagination": pageResult, "message": "success"})
}

func (mc *ManagementController) GetProductModel() {
	var uriParams IDField
	if err := mc.ctx.ShouldBindUri(&uriParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	productModel, err := mc.productModelService.GetProductModel(uriParams.ID)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": productModel, "message": "success"})
}

func (mc *ManagementController) UpdateProductModel() {
	var form models.ProductModel
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	productModel, err := mc.productModelService.GetProductModel(int64(form.ID))
	if err != nil || productModel.ID == 0 {
		mc.ctx.JSON(404, gin.H{"error": "product model not found"})
		return
	}

	productModelMap := utils.StructToMap(form)
	if err := mc.productModelService.UpdateProductModel(productModel, productModelMap); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": form, "message": "success"})
}

func (mc *ManagementController) AddProductionPlan() {
	var form models.ProductionPlan
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.productionPlanService.CreateProductionPlan(&form); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(201, gin.H{"data": form, "message": "success"})
}

func (mc *ManagementController) DeleteProductionPlan() {
	var form IDsField
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.productionPlanService.DeleteProductionPlans(form.IDs); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"message": "success"})
}



func (mc *ManagementController) GetProductionPlans() {
	var queryParams struct{}
	var paginateParams models.PaginationQuery
	if err := mc.ctx.ShouldBindQuery(&queryParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.ctx.ShouldBindQuery(&paginateParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	queryParamsMap := utils.StructToMap(queryParams)
	paginateParamsMap := utils.StructToMap(paginateParams)
	productionPlans, pageResult, err := mc.productionPlanService.GetProductionPlans(queryParamsMap, paginateParamsMap)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": productionPlans, "pagination": pageResult, "message": "success"})
}

func (mc *ManagementController) GetProductionPlan() {
	var uriParams IDField
	if err := mc.ctx.ShouldBindUri(&uriParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	productionPlan, err := mc.productionPlanService.GetProductionPlan(uriParams.ID)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": productionPlan, "message": "success"})
}

func (mc *ManagementController) UpdateProductionPlan() {
	var form models.ProductionPlan
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	productionPlan, err := mc.productionPlanService.GetProductionPlan(int64(form.ID))
	if err != nil || productionPlan.ID == 0 {
		mc.ctx.JSON(404, gin.H{"error": "production plan not found"})
		return
	}

	productionPlanMap := utils.StructToMap(form)
	if err := mc.productionPlanService.UpdateProductionPlan(productionPlan, productionPlanMap); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": form, "message": "success"})
}

func (mc *ManagementController) GetProductionPlansByDateRange() {
	var queryParams struct {
		Date string `form:"date" binding:"required"`
	}
	if err := mc.ctx.ShouldBindQuery(&queryParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 解析日期参数
	baseDate, err := time.Parse("2006-01-02", queryParams.Date)
	if err != nil {
		mc.ctx.JSON(400, gin.H{"error": "invalid date format, expected YYYY-MM-DD"})
		return
	}

	// 调用服务方法获取生产计划
	productionPlans, err := mc.productionPlanService.GetProductionPlansByDateRange(baseDate)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	mc.ctx.JSON(200, gin.H{"data": productionPlans, "message": "success"})
}

func (mc *ManagementController) GetProductLines() {
	var queryParams struct{}
	var paginateParams models.PaginationQuery
	if err := mc.ctx.ShouldBindQuery(&queryParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.ctx.ShouldBindQuery(&paginateParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	queryParamsMap := utils.StructToMap(queryParams)
	paginateParamsMap := utils.StructToMap(paginateParams)
	productLines, pageResult, err := mc.productLineService.GetProductLines(queryParamsMap, paginateParamsMap)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": productLines, "pagination": pageResult, "message": "success"})
}

func (mc *ManagementController) GetProductLine() {
	var uriParams IDField
	if err := mc.ctx.ShouldBindUri(&uriParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	productLine, err := mc.productLineService.GetProductLine(uriParams.ID)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": productLine, "message": "success"})
}

func (mc *ManagementController) AddProductLine() {
	var form models.ProductLine
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.productLineService.CreateProductLine(&form); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(201, gin.H{"data": form, "message": "success"})
}

func (mc *ManagementController) DeleteProductLine() {
	var form IDsField
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.productLineService.DeleteProductLines(form.IDs); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"message": "success"})
}

func (mc *ManagementController) GetPallets() {
	var queryParams struct{}
	var paginateParams models.PaginationQuery
	if err := mc.ctx.ShouldBindQuery(&queryParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.ctx.ShouldBindQuery(&paginateParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	queryParamsMap := utils.StructToMap(queryParams)
	paginateParamsMap := utils.StructToMap(paginateParams)
	pallets, pageResult, err := mc.palletService.GetPallets(queryParamsMap, paginateParamsMap)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": pallets, "pagination": pageResult, "message": "success"})
}

func (mc *ManagementController) GetPallet() {
	var uriParams IDField
	if err := mc.ctx.ShouldBindUri(&uriParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	pallet, err := mc.palletService.GetPallet(uriParams.ID)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": pallet, "message": "success"})
}

func (mc *ManagementController) GetProducts() {
	var queryParams struct {
		StartTime     string `form:"startTime"`     // 开始时间 YYYY-MM-DD HH:MM:SS
		EndTime       string `form:"endTime"`       // 结束时间 YYYY-MM-DD HH:MM:SS
		Search        string `form:"search"`        // 综合模糊查询（托盘SN、产品SN、产品型号SAP）
		ProductLineID uint   `form:"productLineId"` // 产线ID
		PalletID      uint   `form:"palletId"`      // 托盘ID
		HasDefect     bool   `form:"hasDefect"`     // 是否有缺陷
		DefectReason  string `form:"defectReason"`  // 缺陷原因
	}
	var paginateParams models.PaginationQuery
	if err := mc.ctx.ShouldBindQuery(&queryParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.ctx.ShouldBindQuery(&paginateParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 构建动态查询条件
	queryParamsMap := make(map[string]interface{})

	// 时间范围查询
	var sqlHandlers []func(*gorm.DB) *gorm.DB
	if queryParams.StartTime != "" {
		sqlHandlers = append(sqlHandlers, func(db *gorm.DB) *gorm.DB {
			return db.Where("products.created_at >= ?", queryParams.StartTime)
		})
	}
	if queryParams.EndTime != "" {
		sqlHandlers = append(sqlHandlers, func(db *gorm.DB) *gorm.DB {
			return db.Where("products.created_at <= ?", queryParams.EndTime)
		})
	}

	// 综合模糊查询（托盘SN、产品SN、产品型号SAP）
	if queryParams.Search != "" {
		sqlHandlers = append(sqlHandlers, func(db *gorm.DB) *gorm.DB {
			searchPattern := "%" + queryParams.Search + "%"
			return db.Joins("LEFT JOIN pallets ON products.pallet_id = pallets.id").
				Joins("LEFT JOIN product_models ON products.product_model_id = product_models.id").
				Where("products.sn LIKE ? OR pallets.sn LIKE ? OR product_models.sap LIKE ?",
					searchPattern, searchPattern, searchPattern)
		})
	}

	// 产线ID精确查询
	if queryParams.ProductLineID > 0 {
		queryParamsMap["product_line_id"] = queryParams.ProductLineID
	}
	// 托盘ID精确查询
	if queryParams.PalletID > 0 {
		queryParamsMap["pallet_id"] = queryParams.PalletID
	}

	// 是否有缺陷查询
	if queryParams.HasDefect {
		sqlHandlers = append(sqlHandlers, func(db *gorm.DB) *gorm.DB {
			return db.Where("products.defect_reason IS NOT NULL AND products.defect_reason <> ''")
		})
	}
	// 缺陷原因模糊查询
	if queryParams.DefectReason != "" {
		sqlHandlers = append(sqlHandlers, func(db *gorm.DB) *gorm.DB {
			defectPattern := "%" + queryParams.DefectReason + "%"
			return db.Where("products.defect_reason LIKE ?", defectPattern)
		})
	}

	paginateParamsMap := utils.StructToMap(paginateParams)
	products, pageResult, err := mc.productService.GetProducts(queryParamsMap, paginateParamsMap, sqlHandlers...)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": products, "pagination": pageResult, "message": "success"})
}

func (mc *ManagementController) GetProduct() {
	var uriParams IDField
	if err := mc.ctx.ShouldBindUri(&uriParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	product, err := mc.productService.GetProduct(uriParams.ID)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": product, "message": "success"})
}

func (mc *ManagementController) AddApi() {
	var form models.API
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.apiService.CreateAPI(&form); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(201, gin.H{"data": form, "message": "success"})
}

func (mc *ManagementController) DeleteApi() {
	var form IDsField
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.apiService.DeleteAPIs(form.IDs); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"message": "success"})
}

func (mc *ManagementController) GetApis() {
	var queryParams struct{}
	var paginateParams models.PaginationQuery
	if err := mc.ctx.ShouldBindQuery(&queryParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.ctx.ShouldBindQuery(&paginateParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	queryParamsMap := utils.StructToMap(queryParams)
	paginateParamsMap := utils.StructToMap(paginateParams)
	apis, pageResult, err := mc.apiService.GetAPIs(queryParamsMap, paginateParamsMap)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": apis, "pagination": pageResult, "message": "success"})
}

func (mc *ManagementController) GetApi() {
	var uriParams IDField
	if err := mc.ctx.ShouldBindUri(&uriParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	api, err := mc.apiService.GetAPI(uriParams.ID)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": api, "message": "success"})
}

func (mc *ManagementController) UpdateApi() {
	var form models.API
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	api, err := mc.apiService.GetAPI(int64(form.ID))
	if err != nil || api.ID == 0 {
		mc.ctx.JSON(404, gin.H{"error": "api not found"})
		return
	}

	apiMap := utils.StructToMap(form)
	if err := mc.apiService.UpdateAPI(api, apiMap); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": form, "message": "success"})
}

func (mc *ManagementController) AddUser() {
	var form models.User
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.userService.CreateUser(&form); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(201, gin.H{"data": form, "message": "success"})
}

func (mc *ManagementController) DeleteUser() {
	var form IDsField
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.userService.DeleteUsers(form.IDs); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"message": "success"})
}

func (mc *ManagementController) GetUsers() {
	var queryParams struct{}
	var paginateParams models.PaginationQuery
	if err := mc.ctx.ShouldBindQuery(&queryParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := mc.ctx.ShouldBindQuery(&paginateParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	queryParamsMap := utils.StructToMap(queryParams)
	paginateParamsMap := utils.StructToMap(paginateParams)
	users, pageResult, err := mc.userService.GetUsers(queryParamsMap, paginateParamsMap)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": users, "pagination": pageResult, "message": "success"})
}

func (mc *ManagementController) GetUser() {
	var uriParams IDField
	if err := mc.ctx.ShouldBindUri(&uriParams); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := mc.userService.GetUserBy(services.IdentifierTypeID, uriParams.ID)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": user, "message": "success"})
}

func (mc *ManagementController) UpdateUser() {
	var form models.User
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := mc.userService.GetUserBy(services.IdentifierTypeID, int64(form.ID))
	if err != nil || user.ID == 0 {
		mc.ctx.JSON(404, gin.H{"error": "user not found"})
		return
	}

	userMap := utils.StructToMap(form)
	if err := mc.userService.UpdateUser(user, userMap); err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mc.ctx.JSON(200, gin.H{"data": form, "message": "success"})
}

func (mc *ManagementController) Login() {
	var form struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := mc.ctx.ShouldBindJSON(&form); err != nil {
		mc.ctx.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "invalid form",
		})
		return
	}

	userInstance, err := mc.userService.GetUserBy(services.IdentifierTypeEmail, form.Username)
	if err != nil {
		mc.ctx.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "login failed",
		})
		return
	} else {
		if err := userInstance.CheckPassword(form.Password); err != nil {
			mc.ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"message": "login failed",
			})
			return
		}
	}

	mc.ctx.JSON(200, gin.H{
		"message": "login success",
		"token":   mc.jwtService.GenerateToken(form.Username, userInstance.ID, models.JwtServiceRoleAdmin),
	})
}

func (mc *ManagementController) GetQualityStats() {
	var query models.QualityStatsQuery
	if err := mc.ctx.ShouldBindQuery(&query); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", query.StartDate)
	if err != nil {
		mc.ctx.JSON(400, gin.H{"error": "Invalid start date format, expected YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", query.EndDate)
	if err != nil {
		mc.ctx.JSON(400, gin.H{"error": "Invalid end date format, expected YYYY-MM-DD"})
		return
	}

	// 调整结束日期到当天的23:59:59
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// 获取统计数据
	stats, err := mc.qualityStatsService.GetQualityStats(startDate, endDate)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	mc.ctx.JSON(200, gin.H{
		"data":    stats,
		"message": "success",
	})
}

func (mc *ManagementController) GetDefectReport() {
	var query models.DefectReportQuery
	if err := mc.ctx.ShouldBindQuery(&query); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 获取不合格报表数据
	report, err := mc.dataReportService.GetDefectReport(&query)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	mc.ctx.JSON(200, gin.H{
		"data":       report.Items,
		"pagination": report.Pagination,
		"message":    "success",
	})
}

func (mc *ManagementController) GetInspectionReport() {
	var query models.InspectionReportQuery
	if err := mc.ctx.ShouldBindQuery(&query); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 获取检测报表数据
	report, err := mc.dataReportService.GetInspectionReport(&query)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	mc.ctx.JSON(200, gin.H{
		"data":       report.Items,
		"pagination": report.Pagination,
		"message":    "success",
	})
}

func (mc *ManagementController) GetCostReport() {
	var query models.CostReportQuery
	if err := mc.ctx.ShouldBindQuery(&query); err != nil {
		mc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 获取检测费用报表数据
	report, err := mc.dataReportService.GetCostReport(&query)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	mc.ctx.JSON(200, gin.H{
		"data":       report.Items,
		"pagination": report.Pagination,
		"message":    "success",
	})
}

func (mc *ManagementController) ImportProductionPlan() {
	file, err := mc.ctx.FormFile("file")
	if err != nil {
		mc.ctx.JSON(400, gin.H{"error": "file is required"})
		return
	}

	f, err := file.Open()
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": "failed to open file"})
		return
	}
	defer f.Close()

	plans, err := mc.productionPlanService.ImportProductionPlan(f)
	if err != nil {
		mc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	mc.ctx.JSON(200, gin.H{"data": plans, "message": "success"})
}
