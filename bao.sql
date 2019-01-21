USE exchange_db;

-- Create a new table called 'TableName' in schema 'SchemaName'
-- Drop the table if it already exists
DROP TABLE exchange_db.jipiao_exchange;
-- Create the table in the specified schema
CREATE TABLE exchange_db.jipiao_exchange
(
    id  int UNSIGNED  PRIMARY KEY not null auto_increment,
    liaotianbao_id char(30) NOT NULL unique,
    weixin_id char(30) not null,
    self_city char(20) not  null,
    self_arrive char(20) not null,
    self_time   int UNSIGNED not null,
    expect_city char(20) not  null,
    expect_arrive char(20) not  null,
    expect_time  int UNSIGNED not null,
    update_time  int UNSIGNED  not null,
    status tinyint not null DEFAULT 0,
 index(liaotianbao_id, self_city,self_arrive,self_time)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;