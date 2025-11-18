-- 初始化数据库脚本
-- 如果数据库不存在则创建（虽然 docker-compose 已经通过环境变量创建了）
-- 这里可以添加一些初始化数据或配置

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 确保使用 utf8mb4 字符集
ALTER DATABASE hisense_vmi CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;

