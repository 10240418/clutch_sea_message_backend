package models

import "time"

// 不合格报表查询相关结构体
type DefectReportQuery struct {
	StartDate      string `form:"startDate" json:"startDate"`
	EndDate        string `form:"endDate" json:"endDate"`
	SupplierID     *uint  `form:"supplierId" json:"supplierId"`
	ProductModelSN string `form:"productModelSN" json:"productModelSN"`
	PageNum        int    `form:"pageNum" json:"page"`
	PageSize       int    `form:"pageSize" json:"pageSize"` // 移除最大值限制，允许-1表示导出全部
}

type DefectReportItem struct {
	SupplierName   string    `json:"supplierName"`
	QualityDate    time.Time `json:"qualityDate"`
	ProductSN      string    `json:"productSN"`
	ProductModelSN string    `json:"productModelSN"`
	BatchNumber    string    `json:"batchNumber"`
	DefectReason   string    `json:"defectReason"`
	Description    string    `json:"description"`
}

type DefectReportResponse struct {
	Items      []DefectReportItem `json:"items"`
	Pagination PaginationResult   `json:"pagination"`
}

// 检测报表查询相关结构体
type InspectionReportQuery struct {
	ProductModelSN string `form:"productModelSN" json:"productModelSN"` // 物料编码
	BatchNumber    string `form:"batchNumber" json:"batchNumber"`       // 批次号
	SupplierName   string `form:"supplierName" json:"supplierName"`     // 生产厂家
	StartDate      string `form:"startDate" json:"startDate"`           // 开始日期
	EndDate        string `form:"endDate" json:"endDate"`               // 结束日期
	PageNum        int    `form:"pageNum" json:"page"`                  // 页码
	PageSize       int    `form:"pageSize" json:"pageSize"`             // 页大小，-1表示导出全部
}

type InspectionReportItem struct {
	ProductModelSN   string `json:"productModelSN"`   // 物料编码
	BatchNumber      string `json:"batchNumber"`      // 批次号
	InspectionCount  int    `json:"inspectionCount"`  // 检测数量
	QualifiedCount   int    `json:"qualifiedCount"`   // 合格数量
	UnqualifiedCount int    `json:"unqualifiedCount"` // 不合格数量
	SupplierName     string `json:"supplierName"`     // 生产厂家
	InspectionDate   string `json:"inspectionDate"`   // 检测日期(YYYY-MM-DD)
	Description      string `json:"description"`      // 物料描述
	ProductLine      string `json:"productLine"`      // 产线信息
}

type InspectionReportResponse struct {
	Items      []InspectionReportItem `json:"items"`
	Pagination PaginationResult       `json:"pagination"`
}

// 检测费用报表查询相关结构体
type CostReportQuery struct {
	SupplierName   string `form:"supplierName" json:"supplierName"`     // 厂家名称
	ProductModelSN string `form:"productModelSN" json:"productModelSN"` // 物料编码
	MotorType      string `form:"motorType" json:"motorType"`           // 电机类型(ProductModel.Description)
	StartDate      string `form:"startDate" json:"startDate"`           // 开始日期
	EndDate        string `form:"endDate" json:"endDate"`               // 结束日期
	PageNum        int    `form:"pageNum" json:"page"`                  // 页码
	PageSize       int    `form:"pageSize" json:"pageSize"`             // 页大小，-1表示导出全部
}

type CostReportItem struct {
	SupplierName     string `json:"supplierName"`     // 厂家
	ProductModelSN   string `json:"productModelSN"`   // 物料编码
	MotorType        string `json:"motorType"`        // 电机类型
	QualifiedCount   int    `json:"qualifiedCount"`   // 合格数量
	UnqualifiedCount int    `json:"unqualifiedCount"` // 不合格数量
	TotalCount       int    `json:"totalCount"`       // 总数量
	TestDate         string `json:"testDate"`         // 检测日期(YYYY-MM-DD)
}

type CostReportResponse struct {
	Items      []CostReportItem `json:"items"`
	Pagination PaginationResult `json:"pagination"`
}
