CREATE TABLE users (
    id              int primary key not null auto_increment,
    name            varchar(255),
    hash_password   varchar(255),
    email           varchar(255),
    phone           varchar(255),
    created_at      datetime default current_timestamp,
    updated_at      datetime default current_timestamp on update current_timestamp,
    deleted_at      datetime default current_timestamp
);
