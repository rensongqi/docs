- [1 Centos7 安装MySQL 8.0](#1-centos7-安装mysql-80)
- [2 SQL 基础语法](#2-sql-基础语法)
  - [1.1 SQL DML 和 DDL](#11-sql-dml-和-ddl)
  - [1.2 创建表](#12-创建表)
  - [1.3 查询表](#13-查询表)
  - [1.4 WHERE](#14-where)
  - [1.5 ORDER BY](#15-order-by)
  - [1.6 INSERT](#16-insert)
  - [1.7 UPDATE](#17-update)
  - [1.8 DELETE](#18-delete)
- [2 SQL 高级语法](#2-sql-高级语法)
  - [2.1 LIKE](#21-like)
  - [2.2 通配符](#22-通配符)
  - [2.3 IN](#23-in)
  - [2.4 BETWEEN](#24-between)
  - [2.5 Alias](#25-alias)
  - [2.6 Join](#26-join)
  - [2.7 UNION/UNION ALL](#27-unionunion-all)
- [3 SQL 函数](#3-sql-函数)
- [4 SQL 高级操作](#4-sql-高级操作)
  - [4.1 复制表](#41-复制表)
- [5 SQL 内置语法](#5-sql-内置语法)

# 1 Centos7 安装MySQL 8.0

```bash
# 1 如果有老版本，应该先卸载老版本
rpm -qa | grep -i mysql
yum -y remove MySQL-*
rm -rf /etc/my.cnf
rm -rf /root/.mysql_sercret
find / -name mysql | xargs rm -rf

# 2 添加mysql 8.0的源
sudo rpm -Uvh https://dev.mysql.com/get/mysql80-community-release-el7-3.noarch.rpm

# 3 安装mysql 8.0
sudo yum --enablerepo=mysql80-community install mysql-community-server -y

# 4 启动服务并配置开机自启
sudo systemctl start mysqld
sudo systemctl enable mysqld

# 5 查看mysql临时密码
grep "A temporary password" /var/log/mysqld.log

# 6 更改mysql临时密码，默认密码验证很复杂， 常规密码都不能审核
[root@node1 ~]# mysql -uroot -p
mysql> ALTER USER 'root'@'localhost' IDENTIFIED BY '5LaloO!Y';
Query OK, 0 rows affected (0.00 sec)

# 7 查看密码验证策略
mysql> SHOW VARIABLES LIKE "validate_password.%";
+--------------------------------------+--------+
| Variable_name                        | Value  |
+--------------------------------------+--------+
| validate_password.check_user_name    | ON     |
| validate_password.dictionary_file    |        |
| validate_password.length             | 8      |
| validate_password.mixed_case_count   | 1      |
| validate_password.number_count       | 1      |
| validate_password.policy             | MEDIUM |
| validate_password.special_char_count | 1      |
+--------------------------------------+--------+
7 rows in set (0.00 sec)
# 说明
validate_password.length 是密码的最小长度，默认是8，我们把它改成6
输入：set global validate_password.length=6;
validate_password.policy 验证密码的复杂程度，我们把它改成0
输入：set global validate_password.policy=0;
validate_password.check_user_name 用户名检查，用户名和密码不能相同，我们也把它关掉
输入：set global validate_password.check_user_name=off;
# 这样可以设置简单密码：ALTER USER ‘root’@‘localhost’ IDENTIFIED BY ‘12345’;

# 8 配置远程访问
CREATE USER 'root'@'10.66.12.97' IDENTIFIED BY '123456';
GRANT ALL ON *.* TO 'root'@'10.66.12.100' IDENTIFIED BY '123456';
FLUSH Privileges;

# mysql 8.0 默认字符集 CHARSET=utf8mb4
```



# 2 SQL 基础语法

一定要记住，**SQL 对大小写不敏感！**

## 1.1 SQL DML 和 DDL

可以把 SQL 分为两个部分：数据操作语言 (DML) 和 数据定义语言 (DDL)。

SQL (结构化查询语言)是用于执行查询的语法。但是 SQL 语言也包含用于更新、插入和删除记录的语法。

查询和更新指令构成了 SQL 的 DML 部分：

- ***SELECT*** - 从数据库表中获取数据
- ***UPDATE*** - 更新数据库表中的数据
- ***DELETE*** - 从数据库表中删除数据
- ***INSERT INTO*** - 向数据库表中插入数据

SQL 的数据定义语言 (DDL) 部分使我们有能力创建或删除表格。我们也可以定义索引（键），规定表之间的链接，以及施加表间的约束。

SQL 中最重要的 DDL 语句:

- ***CREATE DATABASE*** - 创建新数据库
- ***ALTER DATABASE*** - 修改数据库
- ***CREATE TABLE*** - 创建新表
- ***ALTER TABLE*** - 变更（改变）数据库表
- ***DROP TABLE*** - 删除表
- ***CREATE INDEX*** - 创建索引（搜索键）
- ***DROP INDEX*** - 删除索引

## 1.2 创建表

```sql
CREATE TABLE student
(
ID int primary key NOT NULL AUTO_INCREMENT,
NAME varchar(255),
AGE int
)
```

##  1.3 查询表

```sql
# 1 常用命令
SELECT * FROM student;
SELECT NAME FROM student;
SELECT NAME,AGE student;

# 2 在表中，可能会包含重复值。有时您也许希望仅仅列出不同（distinct）的值，可以使用如下命令
SELECT DISTINCT NAME FROM student;
```

## 1.4 WHERE

**WHERE 子句用于规定选择的标准。**

**语法：**

```sql
SELECT 列名称 FROM 表名称 WHERE 列 运算符 值
```

**示例：**

| 操作符  | 示例                                                         | 描述         |
| :------ | ------------------------------------------------------------ | :----------- |
| =       | `SELECT * FROM student WHERE NAME='RSQ';`                    | 等于         |
| <>      | `SELECT * FROM student WHERE AGE <> 20;`                     | 不等于       |
| >       | `SELECT * FROM student WHERE AGE > 20;`                      | 大于         |
| <       | `SELECT * FROM student WHERE AGE < 20;`                      | 小于         |
| >=      | `SELECT * FROM student WHERE AGE >= 20;`                     | 大于等于     |
| <=      | `SELECT * FROM student WHERE AGE <- 20;`                     | 小于等于     |
| BETWEEN | `SELECT * FROM student WHERE (AGE BETWEEN 20 AND 40);`       | 在某个范围内 |
| LIKE    | `SELECT * FROM student WHERE AGE LIKE 20;`                   | 搜索某种模式 |
| OR      | `SELECT * FROM student WHERE (NAME='RSQ' OR AGE=30);`        | 或           |
| AND     | `SELECT * FROM student WHERE (NAME='RSQ' AND AGE=30);`       | 和           |
| OR AND  | `SELECT * FROM Persons WHERE (NAME='RSQ' OR NAME='RSQ2') AND AGE='30';` |              |

## 1.5 ORDER BY

**ORDER BY 语句用于对结果集进行排序**

```sql
# 1 默认是升序
SELECT * FROM student ORDER BY AGE ;

# 2 降序排序
SELECT * FROM student ORDER BY AGE DESC;

# 3 按照字母和年龄顺序排序
SELECT * FROM student ORDER BY NAME,AGE;

# 4 按照字母和年龄顺序逆序排序
SELECT * FROM student ORDER BY NAME,AGE DESC;
```

## 1.6 INSERT

**INSERT INTO 语句用于向表格中插入新的行。**

**语法：**

```sql
INSERT INTO 表名称 VALUES (值1, 值2,....)
```

**我们也可以指定所要插入数据的列：**

```sql
INSERT INTO table_name (列1, 列2,...) VALUES (值1, 值2,....)
```

**示例：**

```sql
# 不指定列名
INSERT INTO student VALUE (7,"agou",41);

# 指定列名
INSERT INTO student (ID,Name,Age) VALUE (8,"maomoa",18);
```

## 1.7 UPDATE

**Update 语句用于修改表中的数据。**

**语法：**

```sql
UPDATE 表名称 SET 列名称 = 新值 WHERE 列名称 = 某值
```

**示例：**

```sql
# 更新某一行中的一个列
UPDATE student SET NAME='maomao' WHERE NAME='maomoa';

# 更新某一行中的若干列
UPDATE student SET NAME='mao',AGE=21 WHERE NAME='maomao';
```

## 1.8 DELETE

**DELETE 语句用于删除表中的行。**

**语法：**

```sql
DELETE FROM 表名称 WHERE 列名称 = 值
```

**示例：**

```sql
DELETE FROM student WHERE NAME='mao';

# 删除所有行
DELETE FROM student;
DELETE * FROM student;
```



# 2 SQL 高级语法 

**测试数据：**

```sql
CREATE TABLE Persons
(
Id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
LastName varchar(255),
FirstName varchar(255),
Address varchar(255),
City varchar(255)
);

INSERT INTO Persons (LastName,FirstName,Address,City) VALUE ("Adams","John","Oxford Street","London");
INSERT INTO Persons (LastName,FirstName,Address,City) VALUE ("Bush","George","Fifth Avenue","New York");
INSERT INTO Persons (LastName,FirstName,Address,City) VALUE ("Carter","Thomas","Changan Street","Beijing");
```

## 2.1 LIKE

```sql
# 匹配包含a的所有学生
SELECT * FROM Persons WHERE LastName LIKE '%a%';

# 匹配不包含a的学生
SELECT * FROM Persons WHERE LastName NOT LIKE '%a%';
```

## 2.2 通配符

| 通配符                     | 描述                   |
| :------------------------- | :--------------------- |
| %                          | 代表零个或多个字符     |
| _                          | 仅替代一个字符         |
| [charlist]                 | 字符列中的任何单一字符 |
| [^charlist]或者[!charlist] | 不在字符列中的任何单   |

```sql
# % ：选取居住在Ne开头的城市中 
SELECT * FROM Persons WHERE CITY LIKE 'Ne%';

# - : 选取第一个字符之后是eorge的Person
SELECT * FROM Persons WHERE FirstName LIKE '_eorge';

# [] ： 我们希望从上面的 "Persons" 表中选取居住的城市以 "A" 或 "L" 或 "N" 开头的人：
SELECT * FROM Persons WHERE City REGEXP '^[ALN]';

# [^] ：我们希望从上面的 "Persons" 表中选取居住的城市不以 "A" 或 "L" 或 "N" 开头的人：
SELECT * FROM Persons WHERE City REGEXP '^[^ALN]';
```

## 2.3 IN

**IN 操作符允许我们在 WHERE 子句中规定多个值。**

**语法：**

```sql
SELECT column_name(s)
FROM table_name
WHERE column_name IN (value1,value2,...)
```

**示例：**

```sql
# 选择城市在Beijing和New York中的行
SELECT * FROM Persons WHERE City IN ("Beijing","New York");
```

## 2.4 BETWEEN

**BETWEEN ... AND** 会选取介于两个值之间的数据范围。这些值可以是数值、文本或者日期。

**语法：**

```
SELECT column_name(s)
FROM table_name
WHERE column_name
BETWEEN value1 AND value2
```

**示例：**

```sql
# 以字母顺序显示介于 "Adams"（包括）和 "Carter"（包括）之间的人：
SELECT * FROM Persons WHERE LastName BETWEEN 'Adams' AND 'Carter';

# 显示范围之外的人，使用 NOT 操作符
SELECT * FROM Persons WHERE LastName NOT BETWEEN 'Adams' AND 'Carter';
```

## 2.5 Alias

**通过使用 SQL，可以为列名称和表名称指定别名（Alias）。**

**表SQL的语法：**

```
SELECT column_name(s)
FROM table_name
AS alias_name
```

**列的 SQL Alias 语法：**

```
SELECT column_name AS alias_name
FROM table_name
```

**表别名示例：**

```sql
# 有两个表，查询表Persons中LastName为‘Carter’ 然后 Student中NAME为'RSQ'的行
SELECT p.LastName,s.NAME 
FROM Persons AS p,Student AS s 
WHERE p.LastName = 'Carter' AND s.NAME = 'RSQ';
```

**列别名示例：**

```sql
mysql> SELECT LastName AS Family, FirstName AS Name FROM Persons;
+--------+--------+
| Family | Name   |
+--------+--------+
| Adams  | John   |
| Bush   | George |
| Carter | Thomas |
+--------+--------+
```

## 2.6 Join

**SQL join 用于根据两个或多个表中的列之间的关系，从这些表中查询数据。**

**Student表**

```sql
mysql> SELECT * FROM Student;
+----+----------+------+
| ID | NAME     | AGE  |
+----+----------+------+
|  1 | RSQ      |   20 |
|  2 | aust     |   30 |
|  3 | test     |   27 |
|  4 | zhangsan |   38 |
|  5 | lisi     |   28 |
|  6 | wangwu   |   47 |
|  7 | agou     |   41 |
+----+----------+------+
```

**Persons表**

```sql
mysql> SELECT * FROM Persons;
+----+----------+-----------+----------------+----------+
| Id | LastName | FirstName | Address        | City     |
+----+----------+-----------+----------------+----------+
|  1 | Adams    | John      | Oxford Street  | London   |
|  2 | Bush     | George    | Fifth Avenue   | New York |
|  3 | Carter   | Thomas    | Changan Street | Beijing  |
+----+----------+-----------+----------------+----------+
```



几种不同的`JOIN`举例：

- `INNER JOIN`: 如果表中有至少一个匹配，则返回行。`INNER JOIN` 与 `JOIN` 是相同的。

```sql
mysql> SELECT p.LastName, p.FirstName, s.NAME, s.AGE
    -> FROM Persons AS p
    -> INNER JOIN Student AS s
    -> ON p.Id = s.ID
    -> ORDER BY p.LastName;
+----------+-----------+------+------+
| LastName | FirstName | NAME | AGE  |
+----------+-----------+------+------+
| Adams    | John      | RSQ  |   20 |
| Bush     | George    | aust |   30 |
| Carter   | Thomas    | test |   27 |
+----------+-----------+------+------+
```

- `LEFT JOIN`: 即使右表中没有匹配，也从左表返回所有的行

```sql
mysql> SELECT p.LastName, p.FirstName, s.NAME, s.AGE
    -> FROM Persons AS p
    -> LEFT JOIN Student AS s
    -> ON p.Id = s.ID
    -> ORDER BY p.LastName;
+----------+-----------+------+------+
| LastName | FirstName | NAME | AGE  |
+----------+-----------+------+------+
| Adams    | John      | RSQ  |   20 |
| Bush     | George    | aust |   30 |
| Carter   | Thomas    | test |   27 |
+----------+-----------+------+------+
```

- `RIGHT JOIN`: 即使左表中没有匹配，也从右表返回所有的行

```sql
mysql> SELECT p.LastName, p.FirstName, s.NAME, s.AGE
    -> FROM Persons AS p
    -> RIGHT JOIN Student AS s
    -> ON p.Id = s.ID
    -> ORDER BY p.LastName;
+----------+-----------+----------+------+
| LastName | FirstName | NAME     | AGE  |
+----------+-----------+----------+------+
| NULL     | NULL      | zhangsan |   38 |
| NULL     | NULL      | lisi     |   28 |
| NULL     | NULL      | wangwu   |   47 |
| NULL     | NULL      | agou     |   41 |
| Adams    | John      | RSQ      |   20 |
| Bush     | George    | aust     |   30 |
| Carter   | Thomas    | test     |   27 |
+----------+-----------+----------+------+
```

## 2.7 UNION/UNION ALL

UNION 操作符用于合并两个或多个 SELECT 语句的结果集。

UNION 内部的 SELECT 语句必须拥有相同数量的列。列也必须拥有相似的数据类型。同时，每条 SELECT 语句中的列的顺序必须相同。

```sql
mysql> SELECT * FROM Persons UNION SELECT * FROM Data;
+----+----------+-----------+----------------+----------+
| Id | LastName | FirstName | Address        | City     |
+----+----------+-----------+----------------+----------+
|  1 | Adams    | John      | Oxford Street  | London   |
|  2 | Bush     | George    | Fifth Avenue   | New York |
|  3 | Carter   | Thomas    | Changan Street | Beijing  |
+----+----------+-----------+----------------+----------+
```



默认地，UNION 操作符选取不同的值。如果允许重复的值，请使用 UNION ALL。

```sql
mysql> SELECT * FROM Persons UNION ALL SELECT * FROM Data;
+----+----------+-----------+----------------+----------+
| Id | LastName | FirstName | Address        | City     |
+----+----------+-----------+----------------+----------+
|  1 | Adams    | John      | Oxford Street  | London   |
|  2 | Bush     | George    | Fifth Avenue   | New York |
|  3 | Carter   | Thomas    | Changan Street | Beijing  |
|  1 | Adams    | John      | Oxford Street  | London   |
|  2 | Bush     | George    | Fifth Avenue   | New York |
|  3 | Carter   | Thomas    | Changan Street | Beijing  |
+----+----------+-----------+----------------+----------+
```



# 3 SQL 函数





# 4 SQL 高级操作

## 4.1 复制表

```sql
# 1 仅复制表结构
CREATE TABLE edu SELECT * FROM student LIMIT 0;

# 2 复制表数据和表结构
CREATE TABLE edu SELECT * FROM student;

# 3 复制指定列的数据和表结构
CREATE TABLE edu SELECT NAME,AGE FROM student;
```



# 5 SQL 内置语法

```sql
# 1 查看系统变量
SHOW VARIABLES;

# 2 查看当前字符集
mysql> SHOW VARIABLES LIKE '%character%';
+--------------------------+--------------------------------+
| Variable_name            | Value                          |
+--------------------------+--------------------------------+
| character_set_client     | utf8mb4                        |
| character_set_connection | utf8mb4                        |
| character_set_database   | utf8mb4                        |
| character_set_filesystem | binary                         |
| character_set_results    | utf8mb4                        |
| character_set_server     | utf8mb4                        |
| character_set_system     | utf8mb3                        |
| character_sets_dir       | /usr/share/mysql-8.0/charsets/ |
+--------------------------+--------------------------------+

# 3 查看当前连接数
mysql> SHOW VARIABLES LIKE '%max_connect%';
+------------------------+-------+
| Variable_name          | Value |
+------------------------+-------+
| max_connect_errors     | 100   |
| max_connections        | 151   |
| mysqlx_max_connections | 100   |
+------------------------+-------+
```





参考文章：

- [SQL 基础教程 - w3school](https://www.w3school.com.cn/sql/index.asp)
- [MySQL教程 - 菜鸟教程](https://www.runoob.com/mysql/mysql-transaction.html)