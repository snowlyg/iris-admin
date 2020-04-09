use mysql;
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY '123456';
create database test;
use test;
create table user
(
    id int auto_increment primary key,
    username varchar(64) unique not null,
    email varchar(120) unique not null,
    password_hash varchar(128) not null,
    avatar varchar(128) not null
);
insert into user values(1, "mysql","test12345@qq.com","passwd","avaterpath");
insert into user values(2, "lisi","12345test@qq.com","passwd","avaterpath");
