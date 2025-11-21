package controllers

import (
	"fmt"
	"time"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/services"
	"github.com/dreamskynl/godi"
	"github.com/gin-gonic/gin"
)

type ProductionController struct {
	ctx                   *gin.Context
	productLineService    services.IProductLineService
	palletService         services.IPalletService
	productService        services.IProductService
	productModelService   services.IProductModelService
	productionPlanService services.IProductionPlanService
	supplierService       services.ISupplierService
	keyManagementService  services.IKeyManagementService
	jwtService            services.IJwtService
}

func NewProductionController(ctx *gin.Context, sc godi.IGoDI) IProductionController {
	return &ProductionController{
		ctx:                   ctx,
		productLineService:    sc.MustResolve(&services.ProductLineService{}).(*services.ProductLineService),
		palletService:         sc.MustResolve(&services.PalletService{}).(*services.PalletService),
		productService:        sc.MustResolve(&services.ProductService{}).(*services.ProductService),
		productModelService:   sc.MustResolve(&services.ProductModelService{}).(*services.ProductModelService),
		productionPlanService: sc.MustResolve(&services.ProductionPlanService{}).(*services.ProductionPlanService),
		supplierService:       sc.MustResolve(&services.SupplierService{}).(*services.SupplierService),
		keyManagementService:  sc.MustResolve(&services.KeyManagementService{}).(*services.KeyManagementService),
		jwtService:            sc.MustResolve(&services.JwtService{}).(*services.JwtService),
	}
}

func (pc *ProductionController) AddProductLine() {
	var form models.ProductLine
	if err := pc.ctx.ShouldBindJSON(&form); err != nil {
		pc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := pc.productLineService.CreateProductLine(&form); err != nil {
		pc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	pc.ctx.JSON(201, gin.H{"data": form, "message": "success"})
}

func (pc *ProductionController) DeleteProductLine() {
	var form IDsField
	if err := pc.ctx.ShouldBindJSON(&form); err != nil {
		pc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := pc.productLineService.DeleteProductLines(form.IDs); err != nil {
		pc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	pc.ctx.JSON(200, gin.H{"message": "success"})
}

func (pc *ProductionController) AddPallet() {
	var form struct {
		SN   string `json:"sn" binding:"required"`              // Pallet SN
		SAP  string `json:"productModelSap" binding:"required"` // SAP Code
		Goal int    `json:"goal" binding:"required"`
	}
	if err := pc.ctx.ShouldBindJSON(&form); err != nil {
		pc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 从JWT token中获取产线ID
	var productLineID *uint
	if id, exists := pc.ctx.Get("id"); exists {
		fmt.Println(id)
		if lineID, ok := id.(int64); ok {
			// 验证产线是否存在
			_, err := pc.productLineService.GetProductLine(int64(lineID))
			lineIDUint := uint(lineID)
			if err == nil {
				productLineID = &lineIDUint
			}
		}
	}

	// 根据SN查找产品型号
	var productModelID *uint
	productModel, err := pc.productModelService.GetProductModelBySN(form.SAP)
	if err == nil {
		modelID := uint(productModel.ID)
		productModelID = &modelID
	} else {
		// 创建新的产品型号
		newProductModel := models.ProductModel{
			SN:          form.SAP,
			Description: "Auto-created product model for SN: " + form.SAP,
		}

		if err := pc.productModelService.CreateProductModel(&newProductModel); err != nil {
			pc.ctx.JSON(500, gin.H{"error": "failed to create new product model: " + err.Error()})
			return
		}

		modelID := uint(newProductModel.ID)
		productModelID = &modelID
	}

	// 创建托盘
	pallet := models.Pallet{
		SN:             form.SN,
		ProductModelID: productModelID,
		ProductLineID:  productLineID,
		Goal:           form.Goal,
	}

	if err := pc.palletService.CreatePallet(&pallet); err != nil {
		pc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	pc.ctx.JSON(201, gin.H{"data": pallet, "message": "success"})
}

func (pc *ProductionController) AddProduct() {
	var form struct {
		SN           string `json:"sn" binding:"required"`
		PalletID     uint   `json:"palletId"`
		HasDefect    bool   `json:"hasDefect"`
		DefectReason string `json:"defectReason"`
	}
	if err := pc.ctx.ShouldBindJSON(&form); err != nil {
		pc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if form.HasDefect && form.DefectReason == "" {
		pc.ctx.JSON(400, gin.H{"error": "defectReason is required when hasDefect is true"})
		return
	}

	// 提取BatchNumber（第8-11位）
	var batchNumber string
	if len(form.SN) >= 11 {
		batchNumber = form.SN[7:11] // 第8-11位（0-based索引为7-10）
	}

	// 1. 首先根据产品SN前7位，通过GetProductModelBySN函数，进行ProductModel的查询
	var productModelID *uint
	var productLineID *uint
	var productionPlanID *uint

	if len(form.SN) >= 7 {
		snPrefix := form.SN[:7] // 前7位
		productModel, err := pc.productModelService.GetProductModelBySN(snPrefix)
		if err == nil && productModel != nil {
			modelID := uint(productModel.ID)
			productModelID = &modelID

			// 2. 然后，进行可用生产计划查询
			today := time.Now()
			productionPlan, err := pc.productionPlanService.GetActiveProductionPlan(today, productModelID, true)
			if err == nil && productionPlan != nil {
				planID := uint(productionPlan.ID)
				productionPlanID = &planID
			}
		}
	}

	var product models.Product
	// 3. 最后，根据是否不良来区别创建产品记录
	if !form.HasDefect {
		// 对于正常产品，需要获取产线ID和托盘ID
		if form.PalletID != 0 {
			// 根据托盘ID查询托盘信息获取产线ID
			pallet, err := pc.palletService.GetPallet(int64(form.PalletID))
			if err != nil || pallet == nil {
				pc.ctx.JSON(404, gin.H{"error": "pallet not found"})
				return
			}
			productLineID = pallet.ProductLineID
		}

		product = models.Product{
			SN:               form.SN,
			BatchNumber:      batchNumber,
			ProductModelID:   productModelID,
			ProductLineID:    productLineID,
			ProductionPlanID: productionPlanID,
			PalletID:         &form.PalletID,
			HasDefect:        form.HasDefect,
			DefectReason:     "",
		}
	} else {
		// 对于不良品，产线id、托盘id都为nil
		product = models.Product{
			SN:               form.SN,
			BatchNumber:      batchNumber,
			ProductModelID:   productModelID,
			ProductLineID:    nil,
			ProductionPlanID: productionPlanID,
			PalletID:         nil,
			HasDefect:        form.HasDefect,
			DefectReason:     form.DefectReason,
		}
	}

	// 保存产品记录
	if err := pc.productService.CreateProduct(&product); err != nil {
		pc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	pc.ctx.JSON(201, gin.H{"data": product, "message": "success"})
}

func (pc *ProductionController) RegisterProductLine() {
	var form struct {
		DeviceID       string `json:"deviceId" binding:"required"`
		Name           string `json:"name" binding:"required"`
		PalletSnPrefix string `json:"palletSnPrefix" binding:"required"`
	}
	if err := pc.ctx.ShouldBindJSON(&form); err != nil {
		pc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 查询产线是否已经录入且未注册
	var productLines []models.ProductLine
	query := map[string]interface{}{
		"device_id":     form.DeviceID,
		"is_registered": false,
	}
	productLines, _, err := pc.productLineService.GetProductLines(query, map[string]interface{}{})
	if err != nil {
		pc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(productLines) == 0 {
		pc.ctx.JSON(404, gin.H{"error": "device ID not found or already registered"})
		return
	}

	productLine := &productLines[0]

	// 生成公钥
	publicKey, err := pc.keyManagementService.GeneratePublicKeyFromDeviceID(form.DeviceID)
	if err != nil {
		pc.ctx.JSON(500, gin.H{"error": "failed to generate public key"})
		return
	}

	// 更新产线注册状态、公钥和其他信息
	updateData := map[string]interface{}{
		"name":             form.Name,
		"pallet_sn_prefix": form.PalletSnPrefix,
		"is_registered":    true,
		"public_key":       publicKey,
	}
	if err := pc.productLineService.UpdateProductLine(productLine, updateData); err != nil {
		pc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	pc.ctx.JSON(200, gin.H{
		"message":   "registration successful",
		"publicKey": publicKey,
	})
}

func (pc *ProductionController) AuthenticateProductLine() {
	var form struct {
		DeviceID  string `json:"deviceId" binding:"required"`
		PublicKey string `json:"publicKey" binding:"required"`
	}
	if err := pc.ctx.ShouldBindJSON(&form); err != nil {
		pc.ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 验证DeviceID和PublicKey的匹配性
	if !pc.keyManagementService.ValidateDeviceIDAndPublicKey(form.DeviceID, form.PublicKey) {
		pc.ctx.JSON(401, gin.H{"error": "invalid device ID or public key"})
		return
	}

	// 查询产线是否已注册
	var productLines []models.ProductLine
	query := map[string]interface{}{
		"device_id":     form.DeviceID,
		"is_registered": true,
		"public_key":    form.PublicKey,
	}
	productLines, _, err := pc.productLineService.GetProductLines(query, map[string]interface{}{})
	if err != nil {
		pc.ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(productLines) == 0 {
		pc.ctx.JSON(404, gin.H{"error": "product line not found or not registered"})
		return
	}

	productLine := &productLines[0]

	// 生成JWT token
	token := pc.jwtService.GenerateToken(form.DeviceID, productLine.ID, models.JwtServiceRoleProductionLine)

	pc.ctx.JSON(200, gin.H{
		"message": "authentication successful",
		"token":   token,
	})
}
