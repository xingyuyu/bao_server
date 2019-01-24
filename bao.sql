USE exchange_db;

-- Create a new table called 'TableName' in schema 'SchemaName'
-- Drop the table if it already exists
DROP TABLE exchange_db.jipiao_exchange;
DROP TABLE exchange_db.common_exchange;
-- Create the table in the specified schema
CREATE TABLE exchange_db.jipiao_exchange
(
    id  int UNSIGNED  PRIMARY KEY not null auto_increment,
    liaotianbao_id char(30) NOT NULL unique,
    weixin_id char(30) DEFAULT '',
    self_city char(20) not  null,
    self_arrive char(20) not null,
    self_time   int UNSIGNED not null,
    expect_city char(20) DEFAULT '',
    expect_arrive char(20) DEFAULT '',
    expect_time  int UNSIGNED DEFAULT 0,
    update_time  int UNSIGNED  not null,
    status tinyint not null DEFAULT 0,
    index(liaotianbao_id, self_city,self_arrive,self_time)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE exchange_db.common_exchange
(
    id  int UNSIGNED  PRIMARY KEY not null auto_increment,
    liaotianbao_id char(30) NOT NULL unique,
    weixin_id char(30) not null,
    self_attr char(20) not null,
    expect_attr char(20) not null,
    update_time  int UNSIGNED  not null,
    huodong_type tinyint not null DEFAULT 0,
    status tinyint not null DEFAULT 0,
    index(liaotianbao_id, self_attr, expect_attr)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
