package services

import (
	"fmt"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"gorm.io/gorm"
)

type DataReportService struct {
	db *gorm.DB
}

func NewDataReportService(db *gorm.DB) (IDataReportService, error) {
	return &DataReportService{db: db}, nil
}

func (s *DataReportService) GetDefectReport(query *models.DefectReportQuery) (*models.DefectReportResponse, error) {
	var items []models.DefectReportItem

	// 构建查询
	dbQuery := s.db.Table("products p").
		Select(`
			s.name as supplier_name,
			p.created_at as quality_date,
			p.sn as product_sn,
			pm.sn as product_model_sn,
			pm.description as description,
			p.batch_number,
			p.defect_reason
		`).
		Joins("LEFT JOIN product_models pm ON p.product_model_id = pm.id").
		Joins("LEFT JOIN suppliers s ON pm.supplier_id = s.id").
		Where("p.has_defect = ?", true).
		Where("p.defect_reason != ''")

	// 时间筛选
	if query.StartDate != "" {
		dbQuery = dbQuery.Where("DATE(p.created_at) >= ?", query.StartDate)
	}

	if query.EndDate != "" {
		dbQuery = dbQuery.Where("DATE(p.created_at) <= ?", query.EndDate)
	}

	// 厂家ID筛选
	if query.SupplierID != nil {
		dbQuery = dbQuery.Where("pm.supplier_id = ?", *query.SupplierID)
	}

	// 型号SN搜索
	if query.ProductModelSN != "" {
		dbQuery = dbQuery.Where("pm.sn LIKE ?", "%"+query.ProductModelSN+"%")
	}

	// 计算总数
	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count records: %v", err)
	}

	// 分页参数处理
	pageSize := query.PageSize
	page := query.PageNum
	var isExportAll bool

	if pageSize == -1 {
		// 导出全部数据模式
		isExportAll = true
		pageSize = int(total) // 设置为总数量
		page = 1
	} else {
		// 正常分页模式
		if pageSize <= 0 {
			pageSize = 20
		}
		if page <= 0 {
			page = 1
		}
	}

	// 执行查询
	queryBuilder := dbQuery.Order("p.created_at DESC")

	if !isExportAll {
		// 只有在非导出全部模式下才应用分页
		offset := (page - 1) * pageSize
		queryBuilder = queryBuilder.Limit(pageSize).Offset(offset)
	}

	if err := queryBuilder.Scan(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to query defect report: %v", err)
	}

	// 构建分页结果
	pagination := models.PaginationResult{
		Total:    int(total),
		PageNum:  page,
		PageSize: pageSize,
	}

	return &models.DefectReportResponse{
		Items:      items,
		Pagination: pagination,
	}, nil
}

func (s *DataReportService) GetInspectionReport(query *models.InspectionReportQuery) (*models.InspectionReportResponse, error) {
	// 构建基础SQL查询，按物料编码、批次号、检测日期分组统计
	baseSQL := `
		SELECT 
			pm.sn as product_model_sn,
			pm.description as description,
			p.batch_number,
			DATE(p.created_at) as inspection_date,
			s.name as supplier_name,
			pl.name as product_line,
			COUNT(*) as inspection_count,
			SUM(CASE WHEN p.has_defect = false THEN 1 ELSE 0 END) as qualified_count,
			SUM(CASE WHEN p.has_defect = true THEN 1 ELSE 0 END) as unqualified_count
		FROM products p
		LEFT JOIN product_models pm ON p.product_model_id = pm.id
		LEFT JOIN suppliers s ON pm.supplier_id = s.id
		LEFT JOIN product_lines pl ON p.product_line_id = pl.id
		WHERE 1=1`

	var conditions []string
	var args []interface{}

	// 物料编码筛选
	if query.ProductModelSN != "" {
		conditions = append(conditions, "pm.sn LIKE ?")
		args = append(args, "%"+query.ProductModelSN+"%")
	}

	// 批次号筛选
	if query.BatchNumber != "" {
		conditions = append(conditions, "p.batch_number LIKE ?")
		args = append(args, "%"+query.BatchNumber+"%")
	}

	// 生产厂家筛选
	if query.SupplierName != "" {
		conditions = append(conditions, "s.name LIKE ?")
		args = append(args, "%"+query.SupplierName+"%")
	}

	// 时间范围筛选
	if query.StartDate != "" {
		conditions = append(conditions, "DATE(p.created_at) >= ?")
		args = append(args, query.StartDate)
	}

	if query.EndDate != "" {
		conditions = append(conditions, "DATE(p.created_at) <= ?")
		args = append(args, query.EndDate)
	}

	// 构建完整的查询SQL
	for _, condition := range conditions {
		baseSQL += " AND " + condition
	}
	baseSQL += " GROUP BY pm.sn, pm.description, p.batch_number, DATE(p.created_at), s.name, pl.name ORDER BY inspection_date DESC, pm.sn, p.batch_number"

	// 先查询总数
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM (%s) as temp", baseSQL)
	var total int64
	if err := s.db.Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count inspection report records: %v", err)
	}

	// 分页参数处理
	pageSize := query.PageSize
	page := query.PageNum
	var isExportAll bool

	if pageSize == -1 {
		// 导出全部数据模式
		isExportAll = true
		pageSize = int(total)
		page = 1
	} else {
		// 正常分页模式
		if pageSize <= 0 {
			pageSize = 20
		}
		if page <= 0 {
			page = 1
		}
	}

	// 应用分页
	if !isExportAll {
		offset := (page - 1) * pageSize
		baseSQL += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
	}

	// 执行查询
	var items []models.InspectionReportItem
	if err := s.db.Raw(baseSQL, args...).Scan(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to query inspection report: %v", err)
	}

	// 构建分页结果
	pagination := models.PaginationResult{
		Total:    int(total),
		PageNum:  page,
		PageSize: pageSize,
	}

	return &models.InspectionReportResponse{
		Items:      items,
		Pagination: pagination,
	}, nil
}

func (s *DataReportService) GetCostReport(query *models.CostReportQuery) (*models.CostReportResponse, error) {
	// 构建基础SQL查询，按供应商、物料编码、检测日期分组统计
	baseSQL := `
		SELECT 
			s.name as supplier_name,
			pm.sn as product_model_sn,
			pm.description as motor_type,
			DATE(p.created_at) as test_date,
			SUM(CASE WHEN p.has_defect = false THEN 1 ELSE 0 END) as qualified_count,
			SUM(CASE WHEN p.has_defect = true THEN 1 ELSE 0 END) as unqualified_count,
			COUNT(*) as total_count
		FROM products p
		LEFT JOIN product_models pm ON p.product_model_id = pm.id
		LEFT JOIN suppliers s ON pm.supplier_id = s.id
		WHERE 1=1`

	var conditions []string
	var args []interface{}

	// 厂家名称筛选
	if query.SupplierName != "" {
		conditions = append(conditions, "s.name LIKE ?")
		args = append(args, "%"+query.SupplierName+"%")
	}

	// 物料编码筛选
	if query.ProductModelSN != "" {
		conditions = append(conditions, "pm.sn LIKE ?")
		args = append(args, "%"+query.ProductModelSN+"%")
	}

	// 电机类型筛选
	if query.MotorType != "" {
		conditions = append(conditions, "pm.description LIKE ?")
		args = append(args, "%"+query.MotorType+"%")
	}

	// 时间范围筛选
	if query.StartDate != "" {
		conditions = append(conditions, "DATE(p.created_at) >= ?")
		args = append(args, query.StartDate)
	}

	if query.EndDate != "" {
		conditions = append(conditions, "DATE(p.created_at) <= ?")
		args = append(args, query.EndDate)
	}

	// 构建完整的查询SQL
	for _, condition := range conditions {
		baseSQL += " AND " + condition
	}
	baseSQL += " GROUP BY s.name, pm.sn, pm.description, DATE(p.created_at) ORDER BY test_date DESC, s.name, pm.sn"

	// 先查询总数
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM (%s) as temp", baseSQL)
	var total int64
	if err := s.db.Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count cost report records: %v", err)
	}

	// 分页参数处理
	pageSize := query.PageSize
	page := query.PageNum
	var isExportAll bool

	if pageSize == -1 {
		// 导出全部数据模式
		isExportAll = true
		pageSize = int(total)
		page = 1
	} else {
		// 正常分页模式
		if pageSize <= 0 {
			pageSize = 20
		}
		if page <= 0 {
			page = 1
		}
	}

	// 应用分页
	if !isExportAll {
		offset := (page - 1) * pageSize
		baseSQL += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
	}

	// 执行查询
	var items []models.CostReportItem
	if err := s.db.Raw(baseSQL, args...).Scan(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to query cost report: %v", err)
	}

	// 构建分页结果
	pagination := models.PaginationResult{
		Total:    int(total),
		PageNum:  page,
		PageSize: pageSize,
	}

	return &models.CostReportResponse{
		Items:      items,
		Pagination: pagination,
	}, nil
}
