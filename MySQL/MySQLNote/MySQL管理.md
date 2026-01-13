# 系统数据库

MySQL 数据库安装完成后，自带了以下四个数据库，具体作用如下：

| 数据库 | 含义 |
| --- | --- |
| **mysql** | 存储 MySQL 服务器正常运行所需要的各种信息（时区、主从、用户、权限等） |
| **information_schema** | 提供了访问数据库元数据的各种表和视图，包含数据库、表、字段类型及访问权限等 |
| **performance_schema** | 为 MySQL 服务器运行状态提供了一个底层监控功能，主要用于收集数据库服务器性能参数 |
| **sys** | 包含了一系列方便 DBA 和开发人员利用 performance_schema 性能数据库进行性能调优和诊断的视图 |

# 常用工具

## 1. mysqlbinlog

由于服务器生成的二进制日志文件以二进制格式保存，所以如果想要检查这些文本的文本格式，就会使用到 `mysqlbinlog` 日志管理工具。

* **语法：**
`mysqlbinlog [options] log-files1 log-files2 ...`
* **选项：**
* `-d, --database=name`: 指定数据库名称，只列出指定的数据库相关操作。
* `-o, --offset=#`: 忽略掉日志中的前 n 行命令。
* `-r, --result-file=name`: 将输出的文本格式日志输出到指定文件。
* `-s, --short-form`: 显示简单格式，省略掉一些信息。
* `--start-datetime=date1 --stop-datetime=date2`: 指定日期间隔内的所有日志。
* `--start-position=pos1 --stop-position=pos2`: 指定位置间隔内的所有日志。



---

## 2. mysqlshow

`mysqlshow` 客户端对象查找工具，用来很快地查找存在哪些数据库、数据库中的表、表中的列或者索引。

* **语法：**
`mysqlshow [options] [db_name [table_name [col_name]]]`
* **选项：**
* `--count`: 显示数据库及表的统计信息（数据库，表 均可以不指定）。
* `-i`: 显示指定数据库或者指定表的状态信息。


* **示例：**
* `# 查询每个数据库的表的数量及表中记录的数量`
`mysqlshow -uroot -p2143 --count`
* `# 查询 test 库中每个表中的字段数，及行数`
`mysqlshow -uroot -p2143 test --count`
* `# 查询 test 库中 book 表的详细情况`
`mysqlshow -uroot -p2143 test book --count`



---

## 3. mysql

该 `mysql` 不是指 mysql 服务，而是指 mysql 的客户端工具。

* **语法：**
`mysql [options] [database]`
* **选项：**
* `-u, --user=name`: 指定用户名。
* `-p, --password[=name]`: 指定密码。
* `-h, --host=name`: 指定服务器 IP 或域名。
* `-P, --port=port`: 指定连接端口。
* `-e, --execute=name`: 执行 SQL 语句并退出。



`-e` 选项可以在 MySQL 客户端执行 SQL 语句，而不用连接到 MySQL 数据库再执行，对于一些批处理脚本，这种方式尤其方便。

* **示例：**
`mysql -uroot -p123456 db01 -e "select * from stu";`

---

## 4. mysqladmin

`mysqladmin` 是一个执行管理操作的客户端程序。可以用它来检查服务器的配置和当前状态、创建并删除数据库等。

* **查看选项：**
`mysqladmin --help`
* **常用命令：**
* `create databasename`: 创建一个新数据库。
* `drop databasename`: 删除一个数据库及其所有表。
* `ping`: 检查 mysqld 是否存活。
* `processlist`: 显示服务中活跃线程的列表。
* `status`: 给出一个简短的服务状态信息。
* `version`: 获取服务的版本信息。


* **示例：**
* `mysqladmin -uroot -p123456 drop 'test01';`
* `mysqladmin -uroot -p123456 version;`

## 5. mysqlimport / source

### mysqlimport

`mysqlimport` 是客户端数据导入工具，用来导入 `mysqldump` 加 `-T` 参数后导出的文本文件。

* **语法：**
`mysqlimport [options] db_name textfile1 [textfile2...]`
* **示例：**
`mysqlimport -uroot -p2143 test /tmp/city.txt`

### source

如果需要导入 SQL 文件，可以使用 MySQL 中的 `source` 指令：

* **语法：**
`source /root/xxxxx.sql`

---

## 6. mysqldump

`mysqldump` 客户端工具用来备份数据库或在不同数据库之间进行数据迁移。备份内容包含创建表及插入表的 SQL 语句。

* **语法：**
* `mysqldump [options] db_name [tables]`
* `mysqldump [options] --database/-B db1 [db2 db3...]`
* `mysqldump [options] --all-databases/-A`


* **连接选项：**
* `-u, --user=name`: 指定用户名。
* `-p, --password[=name]`: 指定密码。
* `-h, --host=name`: 指定服务器 IP 或域名。
* `-P, --port=#`: 指定连接端口。


* **输出选项：**
* `--add-drop-database`: 在每个数据库创建语句前加上 `drop database` 语句。
* `--add-drop-table`: 在每个表创建语句前加上 `drop table` 语句，默认开启；不开启使用 (`--skip-add-drop-table`)。
* `-n, --no-create-db`: 不包含数据库的创建语句。
* `-t, --no-create-info`: 不包含数据表的创建语句。
* `-d, --no-data`: 不包含数据。
* `-T, --tab=name`: 自动生成两个文件：一个 `.sql` 文件，包含创建表结构的语句；一个 `.txt` 文件，包含数据。

