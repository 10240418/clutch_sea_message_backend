

-- 1. 复合索引：同时优化 JOIN 和日期过滤
CREATE INDEX idx_products_model_created ON products(product_model_id, created_at);

-- 2. 优化 products 表的日期和缺陷查询 (quality_stats 接口)
CREATE INDEX idx_products_created_defect 
ON products(created_at, has_defect, defect_reason(100));

-- 3. 优化 product_models 的供应商关联查询
CREATE INDEX idx_product_models_supplier 
ON product_models(supplier_id);

-- 4.优化产品型号SN的精确查询 (生产控制器高频使用)
-- GetProductModelBySN() 在生产流程中频繁调用
CREATE INDEX idx_product_models_sn 
ON product_models(sn);

-- 5. 优化产品型号描述的模糊查询 (GetProducts 接口频繁使用)
-- 注意: 只对 LIKE 'xxx%' 前缀查询有效，LIKE '%xxx%' 仍需全表扫描
CREATE INDEX idx_product_models_description 
ON product_models(description);

-- 6. 优化托盘SN的模糊查询 (search 参数)
CREATE INDEX idx_pallets_sn 
ON pallets(sn);

-- 7. 如果经常按日期查询生产计划
CREATE INDEX idx_production_plans_plan_date 
ON production_plans(plan_date);

-- 8. 优化产品SN查询
CREATE INDEX idx_products_sn ON products(sn);

