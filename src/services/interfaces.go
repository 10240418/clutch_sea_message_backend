package services

import (
	"mime/multipart"
	"time"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type IUserService interface {
	CreateUser(user *models.User) error
	GetUserBy(identifierType UserIdentifierType, value interface{}) (*models.User, error)
	GetUsers(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.User, models.PaginationResult, error)
	UpdateUser(userObj *models.User, user map[string]interface{}) error
	DeleteUsers(ids []int64) error
}

type IProductService interface {
	CreateProduct(product *models.Product) error
	GetProduct(id int64) (models.Product, error)
	GetProducts(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.Product, models.PaginationResult, error)
	UpdateProduct(productInstance *models.Product, product map[string]interface{}) error
	DeleteProducts(ids []int64) error
}

type ISupplierService interface {
	CreateSupplier(supplier *models.Supplier) error
	GetSupplier(id int64) (*models.Supplier, error)
	GetSuppliers(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.Supplier, models.PaginationResult, error)
	UpdateSupplier(supplierInstance *models.Supplier, supplier map[string]interface{}) error
	DeleteSuppliers(ids []int64) error
}

type IProductModelService interface {
	CreateProductModel(productModel *models.ProductModel) error
	GetProductModel(id int64) (*models.ProductModel, error)
	GetProductModelBySN(sap string) (*models.ProductModel, error)
	GetProductModels(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.ProductModel, models.PaginationResult, error)
	UpdateProductModel(productModelInstance *models.ProductModel, productModel map[string]interface{}) error
	DeleteProductModels(ids []int64) error
}

type IProductionPlanService interface {
	CreateProductionPlan(productionPlan *models.ProductionPlan) error
	GetProductionPlan(id int64) (*models.ProductionPlan, error)
	GetProductionPlans(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.ProductionPlan, models.PaginationResult, error)
	UpdateProductionPlan(productionPlanInstance *models.ProductionPlan, productionPlan map[string]interface{}) error
	DeleteProductionPlans(ids []int64) error
	GetProductionPlansByDateRange(baseDate time.Time) (map[string][]models.ProductionPlan, error)
	GetActiveProductionPlan(date time.Time, productModelID *uint, allowExceed bool) (*models.ProductionPlan, error)
	ImportProductionPlan(file multipart.File) ([]models.ProductionPlan, error)
}

type IProductLineService interface {
	CreateProductLine(productLine *models.ProductLine) error
	GetProductLine(id int64) (*models.ProductLine, error)
	GetProductLineByDeviceID(deviceID string) (*models.ProductLine, error)
	GetProductLines(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.ProductLine, models.PaginationResult, error)
	UpdateProductLine(productLineInstance *models.ProductLine, productLine map[string]interface{}) error
	DeleteProductLines(ids []int64) error
}

type IPalletService interface {
	CreatePallet(pallet *models.Pallet) error
	GetPallet(id int64) (*models.Pallet, error)
	GetPallets(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.Pallet, models.PaginationResult, error)
	UpdatePallet(palletInstance *models.Pallet, pallet map[string]interface{}) error
	DeletePallets(ids []int64) error
}

type IAPIService interface {
	CreateAPI(api *models.API) error
	GetAPI(id int64) (*models.API, error)
	GetAPIs(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.API, models.PaginationResult, error)
	UpdateAPI(apiInstance *models.API, api map[string]interface{}) error
	DeleteAPIs(ids []int64) error
}

type IJwtService interface {
	GenerateToken(identifier string, id int64, role models.JwtServiceRole) string
	ValidateToken(encodedToken string, role models.JwtServiceRole) (*jwt.Token, error)
}

type IKeyManagementService interface {
	GeneratePublicKeyFromDeviceID(deviceID string) (string, error)
	ValidateDeviceIDAndPublicKey(deviceID, publicKey string) bool
}

type IQualityStatsService interface {
	GetQualityStats(startDate, endDate time.Time) (*models.QualityStatsResponse, error)
}

type IDataReportService interface {
	GetDefectReport(query *models.DefectReportQuery) (*models.DefectReportResponse, error)
	GetInspectionReport(query *models.InspectionReportQuery) (*models.InspectionReportResponse, error)
	GetCostReport(query *models.CostReportQuery) (*models.CostReportResponse, error)
}
