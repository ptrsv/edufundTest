CREATE SCHEMA `edufund` ;


create table users (
	id int not null auto_increment primary key,
    username varchar(255),
    email varchar(255),
    created_at datetime,
    updated_at datetime,
    deleted_at datetime,
    confirm_password varchar(255)
);