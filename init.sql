-- Set timezone using offset instead of name
SET GLOBAL time_zone = '-05:00';
SET time_zone = '-05:00';

-- Use DATETIME type which doesn't auto-convert to UTC
CREATE TABLE IF NOT EXISTS time_log (
    id INT AUTO_INCREMENT PRIMARY KEY,
    timestamp DATETIME NOT NULL COMMENT 'Stores literal time without UTC conversion'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Verify timezone setting
SELECT @@global.time_zone, @@session.time_zone;