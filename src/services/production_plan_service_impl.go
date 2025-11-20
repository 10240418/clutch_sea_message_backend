package services

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ProductionPlanService struct {
	db *gorm.DB
}

func NewProductionPlanService(db *gorm.DB) (IProductionPlanService, error) {
	return &ProductionPlanService{db: db}, nil
}

func (s *ProductionPlanService) CreateProductionPlan(productionPlan *models.ProductionPlan) error {
	return s.db.Create(productionPlan).Error
}

func (s *ProductionPlanService) GetProductionPlan(id int64) (*models.ProductionPlan, error) {
	var productionPlan models.ProductionPlan
	err := s.db.First(&productionPlan, id).Error
	return &productionPlan, err
}

func (s *ProductionPlanService) GetProductionPlans(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.ProductionPlan, models.PaginationResult, error) {
	var productionPlans []models.ProductionPlan
	var pagination models.PaginationResult
	var model = s.db.Model(&models.ProductionPlan{})

	for _, handler := range sqlHandler {
		model = handler(model)
	}
	model = model.Where(query)

	model, pagination = utils.DoPagination(model, paginate)
	
	// Order by ID to maintain import order
	model = model.Order("id ASC")

	result := model.Find(&productionPlans)
	if result.Error != nil {
		return []models.ProductionPlan{}, pagination, result.Error
	}

	return productionPlans, pagination, nil
}

func (s *ProductionPlanService) UpdateProductionPlan(productionPlanInstance *models.ProductionPlan, productionPlan map[string]interface{}) error {
	result := s.db.Model(productionPlanInstance).Updates(productionPlan)
	return result.Error
}

func (s *ProductionPlanService) DeleteProductionPlans(ids []int64) error {
	result := s.db.Delete(&models.ProductionPlan{}, ids)
	return result.Error
}

/*
func (s *ProductionPlanService) GetProductionPlansByDateRange(baseDate time.Time) (map[string][]models.ProductionPlan, error) {
	// 计算T, T+1, T+2, T+3的日期
	t := baseDate
	tPlus1 := baseDate.AddDate(0, 0, 1)
	tPlus2 := baseDate.AddDate(0, 0, 2)
	tPlus3 := baseDate.AddDate(0, 0, 3)

	// 查询指定日期范围内的生产计划
	var productionPlans []models.ProductionPlan
	err := s.db.Preload("ProductModel").
		Where("(start_at <= ? AND end_at >= ?) OR (start_at <= ? AND end_at >= ?) OR (start_at <= ? AND end_at >= ?) OR (start_at <= ? AND end_at >= ?)",
			t.Format("2006-01-02"), t.Format("2006-01-02"), tPlus1.Format("2006-01-02"), tPlus1.Format("2006-01-02"), tPlus2.Format("2006-01-02"), tPlus2.Format("2006-01-02"), tPlus3.Format("2006-01-02"), tPlus3.Format("2006-01-02")).Debug().
		Find(&productionPlans).Error

	if err != nil {
		return nil, err
	}

	// 按日期分组生产计划
	result := map[string][]models.ProductionPlan{
		"T":   {},
		"T+1": {},
		"T+2": {},
		"T+3": {},
	}

	for _, plan := range productionPlans {
		// 检查计划是否包含T日期
		if plan.StartAt.Before(t.AddDate(0, 0, 1)) && plan.EndAt.After(t.AddDate(0, 0, -1)) {
			result["T"] = append(result["T"], plan)
		}
		// 检查计划是否包含T+1日期
		if plan.StartAt.Before(tPlus1.AddDate(0, 0, 1)) && plan.EndAt.After(tPlus1.AddDate(0, 0, -1)) {
			result["T+1"] = append(result["T+1"], plan)
		}
		// 检查计划是否包含T+2日期
		if plan.StartAt.Before(tPlus2.AddDate(0, 0, 1)) && plan.EndAt.After(tPlus2.AddDate(0, 0, -1)) {
			result["T+2"] = append(result["T+2"], plan)
		}
		// 检查计划是否包含T+3日期
		if plan.StartAt.Before(tPlus3.AddDate(0, 0, 1)) && plan.EndAt.After(tPlus3.AddDate(0, 0, -1)) {
			result["T+3"] = append(result["T+3"], plan)
		}
	}

	return result, nil
}
*/

func (s *ProductionPlanService) GetProductionPlansByDateRange(baseDate time.Time) (map[string][]models.ProductionPlan, error) {
    // Temporary stub
    return nil, nil
}

func (s *ProductionPlanService) GetActiveProductionPlan(date time.Time, productModelID *uint, allowExceed bool) (*models.ProductionPlan, error) {
	// This method might need adjustment or removal based on new requirements, 
	// but keeping it for now as it might be used elsewhere.
	// Since the model changed, this implementation is likely broken and needs to be updated if used.
	return nil, nil
}

func (s *ProductionPlanService) ImportProductionPlan(file multipart.File) ([]models.ProductionPlan, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil, err
	}
	defer f.Close()

	// Assuming the first sheet is the one we want
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		fmt.Printf("Error getting rows: %v\n", err)
		return nil, err
	}

	fmt.Printf("Total rows read: %d\n", len(rows))

	var plans []models.ProductionPlan
	// For testing: use fixed date instead of today
	targetDate := "2025-08-21"
	fmt.Printf("Target date for filtering: %s\n", targetDate)

	// Skip header row
	for i, row := range rows {
		if i == 0 {
			fmt.Printf("Header row: %v\n", row)
			continue
		}
		if len(row) < 6 { // Ensure enough columns for basic data
			fmt.Printf("Row %d skipped: not enough columns (%d)\n", i, len(row))
			continue
		}

		// Log first few columns for debugging
		fmt.Printf("Row %d: MaterialCode='%s', PartNumber='%s', Type='%s', Manufacturer='%s', PlanDate='%s'\n", 
			i, row[0], row[1], row[2], row[3], row[4])

		// Check if the plan date is today
		// New order:
		// [0] MaterialCode, [1] PartNumber, [2] Type, [3] Manufacturer, [4] PlanDate, [5] ProductionLine...
		planDateStr := row[4]
		
		// Try multiple date formats
		var planDate time.Time
		var err error
		
		// First try standard format: YYYY-MM-DD
		planDate, err = time.Parse("2006-01-02", planDateStr)
		if err != nil {
			// Try YY-MM-DD format (e.g., 25-03-04 for 2025-03-04)
			planDate, err = time.Parse("06-01-02", planDateStr)
			if err != nil {
				fmt.Printf("Row %d skipped: invalid date format '%s' (%v)\n", i, planDateStr, err)
				continue
			}
		}

		if planDate.Format("2006-01-02") != targetDate {
			fmt.Printf("Row %d skipped: date '%s' is not target date (%s)\n", i, planDate.Format("2006-01-02"), targetDate)
			continue
		}

		fmt.Printf("Row %d: MATCHED! Will be imported.\n", i)

		plan := models.ProductionPlan{
			MaterialCode:   row[0],
			PartNumber:     row[1],
			Type:           row[2],
			Manufacturer:   row[3],
			PlanDate:       planDate,
			ProductionLine: row[5],
		}

		// Helper to parse int safely
		parseInt := func(s string) int {
			val, _ := strconv.Atoi(s)
			return val
		}

		// T
		if len(row) > 6 {
			plan.TPlanned = parseInt(row[6])
		}
		if len(row) > 7 {
			plan.TActual = parseInt(row[7])
		}
		plan.TUnfinished = plan.TPlanned - plan.TActual

		// T+1
		if len(row) > 9 {
			plan.T1Planned = parseInt(row[9])
		}
		if len(row) > 10 {
			plan.T1Actual = parseInt(row[10])
		}
		plan.T1Unfinished = plan.T1Planned - plan.T1Actual

		// T+2
		if len(row) > 12 {
			plan.T2Planned = parseInt(row[12])
		}
		if len(row) > 13 {
			plan.T2Actual = parseInt(row[13])
		}
		plan.T2Unfinished = plan.T2Planned - plan.T2Actual

		// T+3
		if len(row) > 15 {
			plan.T3Planned = parseInt(row[15])
		}
		if len(row) > 16 {
			plan.T3Actual = parseInt(row[16])
		}
		plan.T3Unfinished = plan.T3Planned - plan.T3Actual

		// 计算汇总统计
		plan.TotalPlanned = plan.TPlanned + plan.T1Planned + plan.T2Planned + plan.T3Planned
		plan.TotalInspected = plan.TActual + plan.T1Actual + plan.T2Actual + plan.T3Actual
		plan.TotalUnfinished = plan.TUnfinished + plan.T1Unfinished + plan.T2Unfinished + plan.T3Unfinished
		
		// 计算达成率（避免除以0）
		if plan.TotalPlanned > 0 {
			plan.AchievementRate = float64(plan.TotalInspected) / float64(plan.TotalPlanned) * 100
		} else {
			plan.AchievementRate = 0
		}

		// 特殊物料备注（假设在第23列，索引22）
		if len(row) > 22 {
			plan.SpecialNote = row[22]
		}

		plans = append(plans, plan)
	}

	fmt.Printf("Total plans to save: %d\n", len(plans))

	// Transaction to save
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Delete old plans for target date to avoid duplicates
		deleteResult := tx.Where("plan_date = ?", targetDate).Delete(&models.ProductionPlan{})
		if deleteResult.Error != nil {
			fmt.Printf("Error deleting old plans: %v\n", deleteResult.Error)
			return deleteResult.Error
		}
		fmt.Printf("Deleted %d old records for date %s\n", deleteResult.RowsAffected, targetDate)
		
		if len(plans) > 0 {
			if err := tx.Create(&plans).Error; err != nil {
				fmt.Printf("Error saving plans: %v\n", err)
				return err
			}
			fmt.Printf("Successfully saved %d new records\n", len(plans))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return plans, nil
}
