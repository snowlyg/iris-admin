-- use mysql;
-- ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY '123456';
create database custom;
use custom;
source /root/go/src/BeeCustom/database/sql/bee_custom.sql;
