-- ============================================
-- 海信VMI数据库 - 大量产品测试数据生成脚本 (第二部分)
-- ============================================
-- 用途: 生成大量产品和托盘数据用于性能测试
-- 目标: 生成10万+产品记录，覆盖多个型号、多条产线、多种缺陷
-- ============================================

SET NAMES utf8mb4;

DELIMITER $$

-- ============================================
-- 存储过程: 生成托盘数据
-- ============================================
DROP PROCEDURE IF EXISTS generate_pallets$$
CREATE PROCEDURE generate_pallets(
    IN batch_size INT,           -- 每批生成数量
    IN product_model_id INT,     -- 产品型号ID
    IN product_line_id INT,      -- 产线ID
    IN prefix VARCHAR(10)        -- 托盘SN前缀
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE pallet_sn VARCHAR(50);
    DECLARE date_str VARCHAR(8);
    
    SET date_str = DATE_FORMAT(NOW(), '%Y%m%d');
    
    WHILE i < batch_size DO
        SET pallet_sn = CONCAT(prefix, date_str, LPAD(i, 3, '0'));
        
        INSERT INTO pallets (created_at, updated_at, sn, product_model_id, product_line_id)
        VALUES (NOW(), NOW(), pallet_sn, product_model_id, product_line_id);
        
        SET i = i + 1;
    END WHILE;
END$$

-- ============================================
-- 存储过程: 生成产品数据 (带缺陷率控制)
-- ============================================
DROP PROCEDURE IF EXISTS generate_products$$
CREATE PROCEDURE generate_products(
    IN batch_size INT,           -- 每批生成数量
    IN product_model_id INT,     -- 产品型号ID
    IN product_model_sn VARCHAR(20), -- 产品型号SN (用于生成产品SN)
    IN product_line_id INT,      -- 产线ID
    IN pallet_id INT,            -- 托盘ID
    IN production_plan_id INT,   -- 生产计划ID (可为NULL)
    IN defect_rate INT,          -- 缺陷率 (0-100)
    IN start_date DATE           -- 开始日期
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE product_sn VARCHAR(50);
    DECLARE has_defect TINYINT(1);
    DECLARE defect_reason VARCHAR(255);
    DECLARE random_defect INT;
    DECLARE random_num INT;
    DECLARE created_time DATETIME;
    DECLARE date_offset INT;
    
    WHILE i < batch_size DO
        -- 生成产品SN (格式: 型号SN + 随机数)
        SET product_sn = CONCAT(
            product_model_sn, 
            LPAD(FLOOR(RAND() * 1000), 3, '0'),
            CHAR(65 + FLOOR(RAND() * 26)),
            LPAD(FLOOR(RAND() * 10000), 5, '0'),
            CHAR(65 + FLOOR(RAND() * 26)),
            LPAD(FLOOR(RAND() * 1000), 3, '0')
        );
        
        -- 随机分配日期 (在指定日期前后3天内)
        SET date_offset = FLOOR(RAND() * 7) - 3;
        SET created_time = DATE_ADD(
            TIMESTAMP(start_date, 
                     TIME(
                         CONCAT(
                             LPAD(FLOOR(8 + RAND() * 12), 2, '0'), ':',
                             LPAD(FLOOR(RAND() * 60), 2, '0'), ':',
                             LPAD(FLOOR(RAND() * 60), 2, '0')
                         )
                     )
            ),
            INTERVAL date_offset DAY
        );
        
        -- 根据缺陷率决定是否有缺陷
        SET random_num = FLOOR(RAND() * 100);
        IF random_num < defect_rate THEN
            SET has_defect = 1;
            -- 随机选择缺陷类型
            SET random_defect = FLOOR(RAND() * 100);
            IF random_defect < 25 THEN
                SET defect_reason = '轴承噪音';
            ELSEIF random_defect < 45 THEN
                SET defect_reason = '外观不良';
            ELSEIF random_defect < 65 THEN
                SET defect_reason = '端子变形';
            ELSEIF random_defect < 80 THEN
                SET defect_reason = '铭牌不良';
            ELSEIF random_defect < 90 THEN
                SET defect_reason = '线束压伤';
            ELSE
                SET defect_reason = '其他缺陷';
            END IF;
        ELSE
            SET has_defect = 0;
            SET defect_reason = NULL;
        END IF;
        
        -- 插入产品数据
        INSERT INTO products (
            created_at, updated_at, sn, 
            product_model_id, product_line_id, pallet_id, production_plan_id,
            has_defect, defect_reason
        ) VALUES (
            created_time, created_time, product_sn,
            product_model_id, product_line_id, pallet_id, production_plan_id,
            has_defect, defect_reason
        );
        
        SET i = i + 1;
    END WHILE;
END$$

-- ============================================
-- 存储过程: 批量生成完整的生产批次数据
-- ============================================
DROP PROCEDURE IF EXISTS generate_production_batch$$
CREATE PROCEDURE generate_production_batch(
    IN product_count INT,        -- 要生成的产品数量
    IN product_model_id INT,     -- 产品型号ID
    IN product_model_sn VARCHAR(20), -- 产品型号SN
    IN product_line_id INT,      -- 产线ID
    IN line_prefix VARCHAR(10),  -- 产线托盘前缀
    IN defect_rate INT,          -- 缺陷率 (0-100)
    IN batch_date DATE           -- 批次日期
)
BEGIN
    DECLARE pallet_count INT;
    DECLARE products_per_pallet INT DEFAULT 50;  -- 每个托盘50个产品
    DECLARE current_pallet INT DEFAULT 0;
    DECLARE last_pallet_id INT;
    DECLARE remaining_products INT;
    
    -- 计算需要多少个托盘
    SET pallet_count = CEIL(product_count / products_per_pallet);
    
    -- 生成托盘
    CALL generate_pallets(pallet_count, product_model_id, product_line_id, line_prefix);
    
    -- 获取刚生成的第一个托盘ID
    SELECT MAX(id) - pallet_count + 1 INTO last_pallet_id FROM pallets;
    
    -- 为每个托盘生成产品
    SET remaining_products = product_count;
    WHILE current_pallet < pallet_count DO
        IF remaining_products >= products_per_pallet THEN
            CALL generate_products(
                products_per_pallet,
                product_model_id,
                product_model_sn,
                product_line_id,
                last_pallet_id + current_pallet,
                NULL,
                defect_rate,
                batch_date
            );
            SET remaining_products = remaining_products - products_per_pallet;
        ELSE
            CALL generate_products(
                remaining_products,
                product_model_id,
                product_model_sn,
                product_line_id,
                last_pallet_id + current_pallet,
                NULL,
                defect_rate,
                batch_date
            );
            SET remaining_products = 0;
        END IF;
        
        SET current_pallet = current_pallet + 1;
    END WHILE;
END$$

DELIMITER ;

