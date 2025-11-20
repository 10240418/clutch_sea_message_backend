package controllers

type IProductionController interface {
	AddProductLine()
	DeleteProductLine()
	AddPallet()
	AddProduct()
	RegisterProductLine()
	AuthenticateProductLine()
}

type IManagementController interface {
	AddSupplier()
	DeleteSupplier()
	GetSuppliers()
	GetSupplier()
	UpdateSupplier()

	AddProductModel()
	DeleteProductModel()
	GetProductModels()
	GetProductModel()
	UpdateProductModel()

	AddProductionPlan()
	DeleteProductionPlan()
	GetProductionPlans()
	GetProductionPlan()
	UpdateProductionPlan()
	GetProductionPlansByDateRange()
	ImportProductionPlan()

	GetProductLines()
	GetProductLine()
	AddProductLine()
	DeleteProductLine()

	GetPallets()
	GetPallet()

	GetProducts()
	GetProduct()

	AddApi()
	DeleteApi()
	GetApis()
	GetApi()
	UpdateApi()

	AddUser()
	DeleteUser()
	GetUsers()
	GetUser()
	UpdateUser()

	GetQualityStats()

	GetDefectReport()
	GetInspectionReport()
	GetCostReport()

	Login()
}
