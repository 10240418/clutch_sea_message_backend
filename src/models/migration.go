package models

import "gorm.io/gorm"

func Migrate(db *gorm.DB) {
	// 设置表的默认选项：InnoDB 引擎 + utf8mb4 字符集
	tableOptions := "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci"

	// 需要迁移的所有模型
	models := []interface{}{
		&Supplier{},
		&ProductModel{},
		&ProductLine{},
		&Pallet{},
		&ProductionPlan{},
		&Product{},
		&User{},
		&API{},
	}

	// 批量迁移
	for _, model := range models {
		db.Set("gorm:table_options", tableOptions).AutoMigrate(model)
	}
}
