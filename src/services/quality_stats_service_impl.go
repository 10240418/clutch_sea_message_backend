package services

import (
	"fmt"
	"time"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"gorm.io/gorm"
)

// 优化后的统计数据聚合结构
type statsAggregation struct {
	Date         string
	SupplierID   int64
	SupplierName string
	DefectReason string
	TotalCount   int64
	DefectCount  int64
}

type QualityStatsService struct {
	db *gorm.DB
}

func NewQualityStatsService(db *gorm.DB) (IQualityStatsService, error) {
	return &QualityStatsService{db: db}, nil
}

func (s *QualityStatsService) GetQualityStats(startDate, endDate time.Time) (*models.QualityStatsResponse, error) {
	// 使用优化的单次查询获取所有统计数据
	return s.getAllStatsOptimized(startDate, endDate)
}

// 优化版本：使用一次SQL查询获取所有需要的数据
func (s *QualityStatsService) getAllStatsOptimized(startDate, endDate time.Time) (*models.QualityStatsResponse, error) {
	// 1. 使用单次聚合查询获取所有基础数据
	query := `
		SELECT 
			DATE(p.created_at) as date,
			s.id as supplier_id,
			s.name as supplier_name,
			COALESCE(p.defect_reason, '') as defect_reason,
			COUNT(*) as total_count,
			SUM(CASE WHEN p.has_defect = true THEN 1 ELSE 0 END) as defect_count
		FROM products p
		INNER JOIN product_models pm ON p.product_model_id = pm.id
		INNER JOIN suppliers s ON pm.supplier_id = s.id
		WHERE p.created_at BETWEEN ? AND ?
		GROUP BY DATE(p.created_at), s.id, s.name, p.defect_reason
		ORDER BY date, supplier_name
	`

	var aggregations []statsAggregation
	if err := s.db.Raw(query, startDate, endDate).Scan(&aggregations).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch aggregated stats: %w", err)
	}

	// 2. 从聚合数据中构建所有响应部分
	qualityRate := s.buildQualityRate(aggregations)
	defectTypeDistribution := s.buildDefectTypeDistribution(aggregations)
	supplierDefectTrend := s.buildSupplierDefectTrend(aggregations)
	defectTrendByType := s.buildDefectTrendByType(aggregations)

	return &models.QualityStatsResponse{
		QualityRate:            qualityRate,
		DefectTypeDistribution: defectTypeDistribution,
		SupplierDefectTrend:    supplierDefectTrend,
		DefectTrendByType:      defectTrendByType,
	}, nil
}

// 从聚合数据构建合格率统计
func (s *QualityStatsService) buildQualityRate(aggregations []statsAggregation) models.QualityRateStats {
	var totalCount, defectCount int64

	for _, agg := range aggregations {
		totalCount += agg.TotalCount
		defectCount += agg.DefectCount
	}

	qualifiedCount := totalCount - defectCount
	var qualityRate float64
	if totalCount > 0 {
		qualityRate = float64(qualifiedCount) / float64(totalCount) * 100
	}

	return models.QualityRateStats{
		QualifiedCount: int(qualifiedCount),
		DefectCount:    int(defectCount),
		TotalCount:     int(totalCount),
		QualityRate:    qualityRate,
	}
}

// 从聚合数据构建不良类型分布
func (s *QualityStatsService) buildDefectTypeDistribution(aggregations []statsAggregation) []models.DefectTypeItem {
	defectMap := make(map[string]int64)
	var totalDefectCount int64

	for _, agg := range aggregations {
		if agg.DefectReason != "" && agg.DefectCount > 0 {
			defectMap[agg.DefectReason] += agg.DefectCount
			totalDefectCount += agg.DefectCount
		}
	}

	var defectTypes []models.DefectTypeItem
	for defectReason, count := range defectMap {
		rate := float64(count) / float64(totalDefectCount) * 100
		defectTypes = append(defectTypes, models.DefectTypeItem{
			Type:  defectReason,
			Count: int(count),
			Rate:  rate,
		})
	}

	return defectTypes
}

// 从聚合数据构建供应商不良趋势
func (s *QualityStatsService) buildSupplierDefectTrend(aggregations []statsAggregation) []models.SupplierDefectTrend {
	// 使用嵌套 map: supplierName -> date -> {total, defect}
	supplierDailyMap := make(map[string]map[string]*struct {
		total  int64
		defect int64
	})

	for _, agg := range aggregations {
		if _, exists := supplierDailyMap[agg.SupplierName]; !exists {
			supplierDailyMap[agg.SupplierName] = make(map[string]*struct {
				total  int64
				defect int64
			})
		}

		if _, exists := supplierDailyMap[agg.SupplierName][agg.Date]; !exists {
			supplierDailyMap[agg.SupplierName][agg.Date] = &struct {
				total  int64
				defect int64
			}{}
		}

		supplierDailyMap[agg.SupplierName][agg.Date].total += agg.TotalCount
		supplierDailyMap[agg.SupplierName][agg.Date].defect += agg.DefectCount
	}

	// 转换为响应格式
	var supplierTrends []models.SupplierDefectTrend
	for supplierName, dailyMap := range supplierDailyMap {
		var dailyData []models.DailyDefectRate
		for date, counts := range dailyMap {
			var defectRate float64
			if counts.total > 0 {
				defectRate = float64(counts.defect) / float64(counts.total) * 100
			}
			dailyData = append(dailyData, models.DailyDefectRate{
				Date:        date,
				DefectRate:  defectRate,
				TotalCount:  int(counts.total),
				DefectCount: int(counts.defect),
			})
		}

		if len(dailyData) > 0 {
			supplierTrends = append(supplierTrends, models.SupplierDefectTrend{
				SupplierName: supplierName,
				DailyData:    dailyData,
			})
		}
	}

	return supplierTrends
}

// 从聚合数据构建各类型不良趋势
func (s *QualityStatsService) buildDefectTrendByType(aggregations []statsAggregation) models.DefectTrendByType {
	// 定义缺陷类型映射
	defectTypeMap := map[string]string{
		"端子变形": "terminal",
		"铭牌不良": "tag",
		"外观不良": "appearance",
		"轴承噪音": "noise",
	}

	// 使用 map 存储每种类型的每日数据: typeKey -> date -> count
	typeDailyMap := make(map[string]map[string]int64)
	for key := range defectTypeMap {
		typeDailyMap[key] = make(map[string]int64)
	}

	// 聚合数据
	for _, agg := range aggregations {
		if agg.DefectReason == "" || agg.DefectCount == 0 {
			continue
		}

		// 检查缺陷原因包含哪种类型
		for defectType, _ := range defectTypeMap {
			// 使用精确匹配或包含判断
			if agg.DefectReason == defectType {
				typeDailyMap[defectType][agg.Date] += agg.DefectCount
			}
		}
	}

	// 转换为响应格式
	trends := models.DefectTrendByType{
		TerminalData:   s.convertToDefectCountArray(typeDailyMap["端子变形"]),
		TagData:        s.convertToDefectCountArray(typeDailyMap["铭牌不良"]),
		AppearanceData: s.convertToDefectCountArray(typeDailyMap["外观不良"]),
		NoiseData:      s.convertToDefectCountArray(typeDailyMap["轴承噪音"]),
	}

	return trends
}

// 辅助函数：将 map[date]count 转换为数组
func (s *QualityStatsService) convertToDefectCountArray(dailyMap map[string]int64) []models.DailyDefectCount {
	var dailyData []models.DailyDefectCount
	for date, count := range dailyMap {
		dailyData = append(dailyData, models.DailyDefectCount{
			Date:  date,
			Count: int(count),
		})
	}
	return dailyData
}

// 获取合格率统计
func (s *QualityStatsService) getQualityRateStats(startDate, endDate time.Time) (*models.QualityRateStats, error) {
	var totalCount, defectCount int64

	// 统计总数
	if err := s.db.Model(&models.Product{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// 统计不合格数
	if err := s.db.Model(&models.Product{}).
		Where("created_at BETWEEN ? AND ? AND has_defect = ?", startDate, endDate, true).
		Count(&defectCount).Error; err != nil {
		return nil, err
	}

	qualifiedCount := totalCount - defectCount
	var qualityRate float64
	if totalCount > 0 {
		qualityRate = float64(qualifiedCount) / float64(totalCount) * 100
	}

	return &models.QualityRateStats{
		QualifiedCount: int(qualifiedCount),
		DefectCount:    int(defectCount),
		TotalCount:     int(totalCount),
		QualityRate:    qualityRate,
	}, nil
}

// 获取不良类型分布
func (s *QualityStatsService) getDefectTypeDistribution(startDate, endDate time.Time) ([]models.DefectTypeItem, error) {
	var results []struct {
		DefectReason string
		Count        int64
	}

	// 查询不良产品的缺陷原因分布
	if err := s.db.Model(&models.Product{}).
		Select("defect_reason, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ? AND has_defect = ? AND defect_reason != ''", startDate, endDate, true).
		Group("defect_reason").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	// 获取总的不良数量
	var totalDefectCount int64
	for _, result := range results {
		totalDefectCount += result.Count
	}

	// 转换为响应格式
	var defectTypes []models.DefectTypeItem
	for _, result := range results {
		rate := float64(result.Count) / float64(totalDefectCount) * 100
		defectTypes = append(defectTypes, models.DefectTypeItem{
			Type:  result.DefectReason,
			Count: int(result.Count),
			Rate:  rate,
		})
	}

	return defectTypes, nil
}

// 获取供应商不良趋势
func (s *QualityStatsService) getSupplierDefectTrend(startDate, endDate time.Time) ([]models.SupplierDefectTrend, error) {
	// 首先获取所有供应商
	var suppliers []models.Supplier
	if err := s.db.Find(&suppliers).Error; err != nil {
		return nil, err
	}

	var supplierTrends []models.SupplierDefectTrend

	for _, supplier := range suppliers {
		dailyData, err := s.getSupplierDailyDefectData(supplier.ID, startDate, endDate)
		if err != nil {
			return nil, err
		}

		if len(dailyData) > 0 {
			supplierTrends = append(supplierTrends, models.SupplierDefectTrend{
				SupplierName: supplier.Name,
				DailyData:    dailyData,
			})
		}
	}

	return supplierTrends, nil
}

// 获取供应商每日不良率数据
func (s *QualityStatsService) getSupplierDailyDefectData(supplierID int64, startDate, endDate time.Time) ([]models.DailyDefectRate, error) {
	// 按天分组统计每个供应商的产品数量和不良数量
	var results []struct {
		Date        string
		TotalCount  int64
		DefectCount int64
	}

	query := `
		SELECT 
			DATE(p.created_at) as date,
			COUNT(*) as total_count,
			SUM(CASE WHEN p.has_defect = true THEN 1 ELSE 0 END) as defect_count
		FROM products p
		INNER JOIN product_models pm ON p.product_model_id = pm.id
		WHERE pm.supplier_id = ? 
			AND p.created_at BETWEEN ? AND ?
		GROUP BY DATE(p.created_at)
		ORDER BY date
	`

	if err := s.db.Raw(query, supplierID, startDate, endDate).Scan(&results).Error; err != nil {
		return nil, err
	}

	var dailyData []models.DailyDefectRate
	for _, result := range results {
		var defectRate float64
		if result.TotalCount > 0 {
			defectRate = float64(result.DefectCount) / float64(result.TotalCount) * 100
		}

		dailyData = append(dailyData, models.DailyDefectRate{
			Date:        result.Date,
			DefectRate:  defectRate,
			TotalCount:  int(result.TotalCount),
			DefectCount: int(result.DefectCount),
		})
	}

	return dailyData, nil
}

// 获取各类型不良趋势
func (s *QualityStatsService) getDefectTrendByType(startDate, endDate time.Time) (*models.DefectTrendByType, error) {
	defectTypes := map[string]string{
		"terminal":   "端子变形",
		"tag":        "铭牌不良",
		"appearance": "外观不良",
		"noise":      "轴承噪音",
	}

	trends := &models.DefectTrendByType{}

	for key, defectType := range defectTypes {
		dailyData, err := s.getDefectTypeDailyData(defectType, startDate, endDate)
		if err != nil {
			return nil, err
		}

		switch key {
		case "terminal":
			trends.TerminalData = dailyData
		case "tag":
			trends.TagData = dailyData
		case "appearance":
			trends.AppearanceData = dailyData
		case "noise":
			trends.NoiseData = dailyData
		}
	}

	return trends, nil
}

// 获取特定缺陷类型的每日数据
func (s *QualityStatsService) getDefectTypeDailyData(defectType string, startDate, endDate time.Time) ([]models.DailyDefectCount, error) {
	var results []struct {
		Date  string
		Count int64
	}

	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM products 
		WHERE has_defect = true 
			AND defect_reason LIKE ?
			AND created_at BETWEEN ? AND ?
		GROUP BY DATE(created_at)
		ORDER BY date
	`

	searchPattern := fmt.Sprintf("%%%s%%", defectType)
	if err := s.db.Raw(query, searchPattern, startDate, endDate).Scan(&results).Error; err != nil {
		return nil, err
	}

	var dailyData []models.DailyDefectCount
	for _, result := range results {
		dailyData = append(dailyData, models.DailyDefectCount{
			Date:  result.Date,
			Count: int(result.Count),
		})
	}

	return dailyData, nil
}
