-- ============================================
-- 海信VMI数据库测试数据生成脚本
-- ============================================
-- 生成日期: 2025-11-22
-- 用途: 为性能测试和功能测试生成大量真实模拟数据
-- 注意: 执行前请确认是否在测试环境，避免影响生产数据
-- ============================================

-- 设置字符集
SET NAMES utf8mb4;

-- ============================================
-- 第一部分: 供应商数据 (Suppliers)
-- ============================================
-- 当前有3个供应商，补充到10个

INSERT INTO suppliers (id, created_at, updated_at, name, sap, type) VALUES
(3, NOW(), NOW(), '青岛海信电器股份有限公司', '200100', '直接供应'),
(4, NOW(), NOW(), '浙江正泰电器股份有限公司', '200215', '直接供应'),
(5, NOW(), NOW(), '安徽皖南电机股份有限公司', '200320', '直接供应'),
(6, NOW(), NOW(), '深圳市汇川技术股份有限公司', '200445', '直接供应'),
(7, NOW(), NOW(), '宁波韵升股份有限公司', '200567', '贸易商'),
(8, NOW(), NOW(), '上海电气集团股份有限公司', '200688', '直接供应'),
(9, NOW(), NOW(), '广东威灵电机制造有限公司', '200799', '直接供应'),
(10, NOW(), NOW(), '杭州富生电器有限公司', '200812', '贸易商')
ON DUPLICATE KEY UPDATE 
    updated_at = NOW(),
    name = VALUES(name);

-- ============================================
-- 第二部分: 产线数据 (Product Lines)
-- ============================================
-- 当前有2条产线，补充到6条

INSERT INTO product_lines (id, created_at, updated_at, name, pallet_sn_prefix, device_id, is_registered, public_key) VALUES
(5, NOW(), NOW(), '交流电机全检C线', 'JLC', CONCAT('test_device_', UUID()), 1, SHA2(CONCAT('key_', UUID()), 256)),
(6, NOW(), NOW(), '直流电机全检B线', 'ZLB', CONCAT('test_device_', UUID()), 1, SHA2(CONCAT('key_', UUID()), 256)),
(7, NOW(), NOW(), '交流电机抽检A线', 'JLCA', CONCAT('test_device_', UUID()), 1, SHA2(CONCAT('key_', UUID()), 256)),
(8, NOW(), NOW(), '直流电机抽检A线', 'ZLCA', CONCAT('test_device_', UUID()), 1, SHA2(CONCAT('key_', UUID()), 256))
ON DUPLICATE KEY UPDATE 
    updated_at = NOW(),
    name = VALUES(name);

-- ============================================
-- 第三部分: 产品型号数据 (Product Models)
-- ============================================
-- 当前有8个型号，补充到50个

INSERT INTO product_models (id, created_at, updated_at, sn, part_number, description, supplier_id) VALUES
-- 交流电机系列 (供应商: 14-常州永安, 3-海信, 4-正泰, 8-上海电气)
(2006, NOW(), NOW(), '1068456', NULL, 'H7B02656A/交流电机/YSK-30-4G/220V', 14),
(2007, NOW(), NOW(), '1068457', NULL, 'H7B02657B/交流电机/YSK-35-4G/220V', 3),
(2008, NOW(), NOW(), '1068458', NULL, 'H7B02658C/交流电机/YSK-40-4G/380V', 4),
(2009, NOW(), NOW(), '1068459', NULL, 'H7B02659A/交流电机/YSK-45-6G/380V', 8),
(2010, NOW(), NOW(), '1109238', NULL, 'H7B1109238/交流电机/YDK-50-8G/380V', 14),
(2011, NOW(), NOW(), '1109239', NULL, 'H7B1109239/交流电机/YDK-60-8G/380V', 3),
(2012, NOW(), NOW(), '1109240', NULL, 'H7B1109240/交流电机/YDK-70-8G/380V', 4),

-- 直流电机系列 (供应商: 13-卧龙, 2-上骐, 5-皖南, 6-汇川, 9-威灵)
(2013, NOW(), NOW(), '1086134', NULL, 'H7B04861B-D/直流电机/ZWF-42F-1/310V/否', 13),
(2014, NOW(), NOW(), '1086135', NULL, 'H7B04862C-D/直流电机/ZWF-45F-1/310V/否', 2),
(2015, NOW(), NOW(), '1086136', NULL, 'H7B04863A-D/直流电机/ZWF-48F-2/310V/否', 5),
(2016, NOW(), NOW(), '1203992', NULL, 'H7B10835C-D/直流电机/WZ-65-8A-D2/DC 310V/否', 6),
(2017, NOW(), NOW(), '1203993', NULL, 'H7B10836A-D/直流电机/WZ-70-8A-D2/DC 310V/否', 9),
(2018, NOW(), NOW(), '1203994', NULL, 'H7B10837B-D/直流电机/WZ-75-8A-D2/DC 310V/否', 13),
(2019, NOW(), NOW(), '1203995', NULL, 'H7B10838C-D/直流电机/WZ-80-8A-D2/DC 310V/否', 2),

-- 更多型号
(2020, NOW(), NOW(), '1195296', NULL, 'H7B1195296/直流电机/WZ-55-6A/DC 280V/否', 13),
(2021, NOW(), NOW(), '1195297', NULL, 'H7B1195297/直流电机/WZ-58-6A/DC 280V/否', 2),
(2022, NOW(), NOW(), '1195298', NULL, 'H7B1195298/直流电机/WZ-62-6A/DC 280V/否', 5),
(2023, NOW(), NOW(), '1199730', NULL, 'H7B1199730/直流电机/WZ-50-4A/DC 250V/否', 6),
(2024, NOW(), NOW(), '1199731', NULL, 'H7B1199731/直流电机/WZ-52-4A/DC 250V/否', 9),
(2025, NOW(), NOW(), '1199732', NULL, 'H7B1199732/直流电机/WZ-54-4A/DC 250V/否', 13),

-- 特殊规格
(2026, NOW(), NOW(), '1204000', NULL, 'H7B1204000/直流电机/WZ-100-10A/DC 350V/否', 2),
(2027, NOW(), NOW(), '1204001', NULL, 'H7B1204001/直流电机/WZ-120-12A/DC 350V/否', 13),
(2028, NOW(), NOW(), '1204002', NULL, 'H7B1204002/交流电机/YDK-80-10G/480V', 14),
(2029, NOW(), NOW(), '1204003', NULL, 'H7B1204003/交流电机/YDK-90-10G/480V', 3),
(2030, NOW(), NOW(), '1204004', NULL, 'H7B1204004/交流电机/YDK-100-12G/480V', 4),

-- 高端系列
(2031, NOW(), NOW(), '1204005', NULL, 'H7B1204005/直流电机/WZ-150-15A/DC 400V/是', 6),
(2032, NOW(), NOW(), '1204006', NULL, 'H7B1204006/直流电机/WZ-180-18A/DC 400V/是', 9),
(2033, NOW(), NOW(), '1204007', NULL, 'H7B1204007/交流电机/YDK-120-15G/600V', 8),
(2034, NOW(), NOW(), '1204008', NULL, 'H7B1204008/交流电机/YDK-150-18G/600V', 14),
(2035, NOW(), NOW(), '1204009', NULL, 'H7B1204009/交流电机/YDK-180-20G/600V', 3),

-- 节能系列
(2036, NOW(), NOW(), '1204010', NULL, 'H7B1204010/直流电机/WZ-45-4A-ECO/DC 250V/否', 13),
(2037, NOW(), NOW(), '1204011', NULL, 'H7B1204011/直流电机/WZ-50-4A-ECO/DC 250V/否', 2),
(2038, NOW(), NOW(), '1204012', NULL, 'H7B1204012/交流电机/YDK-40-6G-ECO/380V', 4),
(2039, NOW(), NOW(), '1204013', NULL, 'H7B1204013/交流电机/YDK-45-6G-ECO/380V', 8),
(2040, NOW(), NOW(), '1204014', NULL, 'H7B1204014/交流电机/YDK-50-8G-ECO/380V', 14),

-- 工业级系列
(2041, NOW(), NOW(), '1204015', NULL, 'H7B1204015/直流电机/WZ-200-20A-IND/DC 450V/是', 6),
(2042, NOW(), NOW(), '1204016', NULL, 'H7B1204016/直流电机/WZ-220-22A-IND/DC 450V/是', 9),
(2043, NOW(), NOW(), '1204017', NULL, 'H7B1204017/交流电机/YDK-200-22G-IND/660V', 8),
(2044, NOW(), NOW(), '1204018', NULL, 'H7B1204018/交流电机/YDK-220-25G-IND/660V', 3),
(2045, NOW(), NOW(), '1204019', NULL, 'H7B1204019/交流电机/YDK-250-28G-IND/660V', 14),

-- 特种应用
(2046, NOW(), NOW(), '1204020', NULL, 'H7B1204020/直流电机/WZ-35-2A-SPE/DC 200V/否', 13),
(2047, NOW(), NOW(), '1204021', NULL, 'H7B1204021/直流电机/WZ-38-2A-SPE/DC 200V/否', 2),
(2048, NOW(), NOW(), '1204022', NULL, 'H7B1204022/交流电机/YDK-30-4G-SPE/220V', 4),
(2049, NOW(), NOW(), '1204023', NULL, 'H7B1204023/交流电机/YDK-32-4G-SPE/220V', 8),
(2050, NOW(), NOW(), '1204024', NULL, 'H7B1204024/交流电机/YDK-35-4G-SPE/220V', 14)
ON DUPLICATE KEY UPDATE 
    updated_at = NOW(),
    description = VALUES(description);

-- ============================================
-- 第四部分: 生产计划数据 (Production Plans)
-- ============================================
-- 生成未来7天的生产计划（每天10-15个型号）

-- 2025-11-23 计划
INSERT INTO production_plans (created_at, updated_at, material_code, part_number, type, manufacturer, plan_date, production_line, t_planned, t_actual) VALUES
(NOW(), NOW(), 'MAT-H7B02656A', 'H7B02656A', '交流', '常州永安', '2025-11-23', '线A', 200, 0),
(NOW(), NOW(), 'MAT-H7B02657B', 'H7B02657B', '交流', '海信', '2025-11-23', '线B', 180, 0),
(NOW(), NOW(), 'MAT-H7B04861B-D', 'H7B04861B-D', '直流', '卧龙', '2025-11-23', '线C', 150, 0),
(NOW(), NOW(), 'MAT-H7B04862C-D', 'H7B04862C-D', '直流', '上骐', '2025-11-23', '线A', 160, 0),
(NOW(), NOW(), 'MAT-H7B10835C-D', 'H7B10835C-D', '直流', '汇川', '2025-11-23', '线B', 140, 0);

-- 2025-11-24 计划
INSERT INTO production_plans (created_at, updated_at, material_code, part_number, type, manufacturer, plan_date, production_line, t_planned, t_actual) VALUES
(NOW(), NOW(), 'MAT-H7B02658C', 'H7B02658C', '交流', '正泰', '2025-11-24', '线A', 220, 0),
(NOW(), NOW(), 'MAT-H7B02659A', 'H7B02659A', '交流', '上海电气', '2025-11-24', '线B', 200, 0),
(NOW(), NOW(), 'MAT-H7B10836A-D', 'H7B10836A-D', '直流', '威灵', '2025-11-24', '线C', 170, 0),
(NOW(), NOW(), 'MAT-H7B10837B-D', 'H7B10837B-D', '直流', '卧龙', '2025-11-24', '线A', 180, 0),
(NOW(), NOW(), 'MAT-H7B1195296', 'H7B1195296', '直流', '卧龙', '2025-11-24', '线B', 150, 0);

-- 2025-11-25 计划
INSERT INTO production_plans (created_at, updated_at, material_code, part_number, type, manufacturer, plan_date, production_line, t_planned, t_actual) VALUES
(NOW(), NOW(), 'MAT-H7B1109238', 'H7B1109238', '交流', '常州永安', '2025-11-25', '线A', 190, 0),
(NOW(), NOW(), 'MAT-H7B1109239', 'H7B1109239', '交流', '海信', '2025-11-25', '线B', 210, 0),
(NOW(), NOW(), 'MAT-H7B1195297', 'H7B1195297', '直流', '上骐', '2025-11-25', '线C', 160, 0),
(NOW(), NOW(), 'MAT-H7B1199730', 'H7B1199730', '直流', '汇川', '2025-11-25', '线A', 140, 0),
(NOW(), NOW(), 'MAT-H7B1204000', 'H7B1204000', '直流', '上骐', '2025-11-25', '线B', 100, 0);

-- 2025-11-26 计划
INSERT INTO production_plans (created_at, updated_at, material_code, part_number, type, manufacturer, plan_date, production_line, t_planned, t_actual) VALUES
(NOW(), NOW(), 'MAT-H7B1109240', 'H7B1109240', '交流', '正泰', '2025-11-26', '线A', 200, 0),
(NOW(), NOW(), 'MAT-H7B1204002', 'H7B1204002', '交流', '常州永安', '2025-11-26', '线B', 180, 0),
(NOW(), NOW(), 'MAT-H7B1204001', 'H7B1204001', '直流', '卧龙', '2025-11-26', '线C', 120, 0),
(NOW(), NOW(), 'MAT-H7B1199731', 'H7B1199731', '直流', '威灵', '2025-11-26', '线A', 150, 0),
(NOW(), NOW(), 'MAT-H7B1204005', 'H7B1204005', '直流', '汇川', '2025-11-26', '线B', 90, 0);

-- 2025-11-27 计划
INSERT INTO production_plans (created_at, updated_at, material_code, part_number, type, manufacturer, plan_date, production_line, t_planned, t_actual) VALUES
(NOW(), NOW(), 'MAT-H7B1204003', 'H7B1204003', '交流', '海信', '2025-11-27', '线A', 220, 0),
(NOW(), NOW(), 'MAT-H7B1204004', 'H7B1204004', '交流', '正泰', '2025-11-27', '线B', 200, 0),
(NOW(), NOW(), 'MAT-H7B1204006', 'H7B1204006', '直流', '威灵', '2025-11-27', '线C', 100, 0),
(NOW(), NOW(), 'MAT-H7B1204010', 'H7B1204010', '直流', '卧龙', '2025-11-27', '线A', 170, 0),
(NOW(), NOW(), 'MAT-H7B1204011', 'H7B1204011', '直流', '上骐', '2025-11-27', '线B', 160, 0);

-- 2025-11-28 计划
INSERT INTO production_plans (created_at, updated_at, material_code, part_number, type, manufacturer, plan_date, production_line, t_planned, t_actual) VALUES
(NOW(), NOW(), 'MAT-H7B1204007', 'H7B1204007', '交流', '上海电气', '2025-11-28', '线A', 180, 0),
(NOW(), NOW(), 'MAT-H7B1204008', 'H7B1204008', '交流', '常州永安', '2025-11-28', '线B', 190, 0),
(NOW(), NOW(), 'MAT-H7B1204015', 'H7B1204015', '直流', '汇川', '2025-11-28', '线C', 80, 0),
(NOW(), NOW(), 'MAT-H7B1204016', 'H7B1204016', '直流', '威灵', '2025-11-28', '线A', 90, 0),
(NOW(), NOW(), 'MAT-H7B1204020', 'H7B1204020', '直流', '卧龙', '2025-11-28', '线B', 200, 0);

-- 2025-11-29 计划
INSERT INTO production_plans (created_at, updated_at, material_code, part_number, type, manufacturer, plan_date, production_line, t_planned, t_actual) VALUES
(NOW(), NOW(), 'MAT-H7B1204009', 'H7B1204009', '交流', '海信', '2025-11-29', '线A', 210, 0),
(NOW(), NOW(), 'MAT-H7B1204012', 'H7B1204012', '交流', '正泰', '2025-11-29', '线B', 190, 0),
(NOW(), NOW(), 'MAT-H7B1204017', 'H7B1204017', '交流', '上海电气', '2025-11-29', '线C', 160, 0),
(NOW(), NOW(), 'MAT-H7B1204021', 'H7B1204021', '直流', '上骐', '2025-11-29', '线A', 180, 0),
(NOW(), NOW(), 'MAT-H7B1204022', 'H7B1204022', '交流', '正泰', '2025-11-29', '线B', 170, 0);

-- ============================================
-- 脚本完成提示
-- ============================================
-- 完成第一批测试数据生成
-- 下一步: 继续生成大量产品数据（products）和托盘数据（pallets）
