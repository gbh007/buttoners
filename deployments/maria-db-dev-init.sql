CREATE USER IF NOT EXISTS root@localhost IDENTIFIED BY '12345';
SET PASSWORD FOR root@localhost = PASSWORD('12345');
GRANT ALL ON *.* TO root@localhost WITH GRANT OPTION;
CREATE USER IF NOT EXISTS root@'%' IDENTIFIED BY '12345';
SET PASSWORD FOR root@'%' = PASSWORD('12345');
GRANT ALL ON *.* TO root@'%' WITH GRANT OPTION;

CREATE USER IF NOT EXISTS auth_user@'%' IDENTIFIED BY 'auth_pwd';
SET PASSWORD FOR auth_user@'%' = PASSWORD('auth_pwd');

CREATE DATABASE IF NOT EXISTS auth_db;
GRANT ALL ON auth_db.* TO auth_user@'%'; 

CREATE USER IF NOT EXISTS notification_user@'%' IDENTIFIED BY 'notification_pwd';
SET PASSWORD FOR notification_user@'%' = PASSWORD('notification_pwd');

CREATE DATABASE IF NOT EXISTS notification_db;
GRANT ALL ON notification_db.* TO notification_user@'%'; 
