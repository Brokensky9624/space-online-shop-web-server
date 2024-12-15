以下是 MySQL 中一些常用的 SQL 語法，包括 `DROP TABLE` 等操作的示例。

---

### 1. **資料庫操作**

#### 創建資料庫
```sql
CREATE DATABASE my_database;
```

#### 切換到指定資料庫
```sql
USE my_database;
```

#### 刪除資料庫
```sql
DROP DATABASE my_database;
```

---

### 2. **表操作**

#### 創建資料表
```sql
CREATE TABLE product (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    title VARCHAR(256),
    description TEXT,
    price DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### 查看資料表結構
```sql
DESCRIBE product;
```

#### 查看資料表的所有資訊
```sql
SHOW CREATE TABLE product;
```

#### 修改資料表
- 添加欄位：
```sql
ALTER TABLE product ADD COLUMN stock INT DEFAULT 0;
```

- 修改欄位名稱或資料型別：
```sql
ALTER TABLE product CHANGE COLUMN description desc TEXT;
```

- 刪除欄位：
```sql
ALTER TABLE product DROP COLUMN stock;
```

#### 刪除資料表
```sql
DROP TABLE product;
```

---

### 3. **資料操作**

#### 插入資料
```sql
INSERT INTO product (name, title, description, price)
VALUES ('Laptop', 'Gaming Laptop', 'A powerful gaming laptop', 1200.00);
```

#### 更新資料
```sql
UPDATE product
SET price = 1100.00
WHERE name = 'Laptop';
```

#### 刪除資料
```sql
DELETE FROM product
WHERE name = 'Laptop';
```

#### 查詢資料
- 查詢所有資料：
```sql
SELECT * FROM product;
```

- 條件查詢：
```sql
SELECT * FROM product WHERE price > 1000;
```

- 排序：
```sql
SELECT * FROM product ORDER BY price DESC;
```

- 分頁：
```sql
SELECT * FROM product LIMIT 10 OFFSET 20;
```

---

### 4. **索引操作**

#### 創建索引
- 創建單一欄位索引：
```sql
CREATE INDEX idx_name ON product (name);
```

- 創建多欄位組合索引：
```sql
CREATE INDEX idx_name_price ON product (name, price);
```

#### 刪除索引
```sql
DROP INDEX idx_name ON product;
```

---

### 5. **聯合查詢**

#### INNER JOIN（內聯結）
```sql
SELECT p.name, c.name AS category_name
FROM product p
INNER JOIN category c ON p.category_id = c.id;
```

#### LEFT JOIN（左聯結）
```sql
SELECT p.name, c.name AS category_name
FROM product p
LEFT JOIN category c ON p.category_id = c.id;
```

#### RIGHT JOIN（右聯結）
```sql
SELECT p.name, c.name AS category_name
FROM product p
RIGHT JOIN category c ON p.category_id = c.id;
```

---

### 6. **其他常用語法**

#### 設置自動遞增值
```sql
ALTER TABLE product AUTO_INCREMENT = 1000;
```

#### 計算記錄數量
```sql
SELECT COUNT(*) FROM product;
```

#### 分組與聚合
- 分組查詢：
```sql
SELECT category_id, COUNT(*) AS product_count
FROM product
GROUP BY category_id;
```

- 聚合函數：
```sql
SELECT AVG(price) AS average_price FROM product;
```

#### 查詢重複值
```sql
SELECT name, COUNT(*) AS count
FROM product
GROUP BY name
HAVING count > 1;
```

---

### 7. **用戶與權限操作**

#### 創建用戶
```sql
CREATE USER 'username'@'localhost' IDENTIFIED BY 'password';
```

#### 賦予權限
```sql
GRANT ALL PRIVILEGES ON my_database.* TO 'username'@'localhost';
```

#### 撤銷權限
```sql
REVOKE ALL PRIVILEGES ON my_database.* FROM 'username'@'localhost';
```

#### 刪除用戶
```sql
DROP USER 'username'@'localhost';
```

---

### 總結
這些是 MySQL 中常用的語法操作，包含了資料庫管理、資料表操作、資料查詢與操作、索引、聯合查詢以及用戶權限管理。根據需求選擇相應的語法操作即可！