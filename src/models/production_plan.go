package models

import "time"

// ProductionPlan 对应 'ProductionPlan' 表
// ProductionPlan 对应 'ProductionPlan' 表
type ProductionPlan struct {
	ModelFields    `s2m:"-"`
	MaterialCode   string    `gorm:"type:varchar(50)" json:"materialCode"`   // 物料编码
	PartNumber     string    `gorm:"type:varchar(50)" json:"partNumber"`     // 部品号
	Type           string    `gorm:"type:varchar(20)" json:"type"`           // 直流/交流
	Manufacturer   string    `gorm:"type:varchar(100)" json:"manufacturer"`  // 厂家
	PlanDate       time.Time `gorm:"type:date" json:"planDate"`              // 计划输入日期
	ProductionLine string    `gorm:"type:varchar(50)" json:"productionLine"` // 生产线体

	// T (当天)
	TPlanned    int `gorm:"type:int" json:"tPlanned"`       // T计划数
	TActual     int `gorm:"type:int" json:"tActual"`        // T完成数
	TUnfinished int `gorm:"type:int" json:"tUnfinished"`    // T未完成数

	// T+1
	T1Planned    int `gorm:"type:int" json:"t1Planned"`     // T1计划数
	T1Actual     int `gorm:"type:int" json:"t1Actual"`      // T1完成数
	T1Unfinished int `gorm:"type:int" json:"t1Unfinished"`  // T1未完成数

	// T+2
	T2Planned    int `gorm:"type:int" json:"t2Planned"`     // T2计划数
	T2Actual     int `gorm:"type:int" json:"t2Actual"`      // T2完成数
	T2Unfinished int `gorm:"type:int" json:"t2Unfinished"`  // T2未完成数

	// T+3
	T3Planned    int `gorm:"type:int" json:"t3Planned"`     // T3计划数
	T3Actual     int `gorm:"type:int" json:"t3Actual"`      // T3完成数
	T3Unfinished int `gorm:"type:int" json:"t3Unfinished"`  // T3未完成数

	// 汇总统计
	TotalPlanned    int     `gorm:"type:int" json:"totalPlanned"`       // 计划数（T+T1+T2+T3）
	TotalInspected  int     `gorm:"type:int" json:"totalInspected"`     // 检验数（T+T1+T2+T3完成数）
	TotalUnfinished int     `gorm:"type:int" json:"totalUnfinished"`    // 未完成数（T+T1+T2+T3）
	AchievementRate float64 `gorm:"type:decimal(5,2)" json:"achievementRate"` // 达成率（%）
	SpecialNote     string  `gorm:"type:varchar(200)" json:"specialNote"`     // 特殊物料备注
}
