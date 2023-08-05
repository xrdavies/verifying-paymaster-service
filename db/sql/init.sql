create database paymaster;
create user paymaster with encrypted password 'paymaster';
grant all privileges on database paymaster to paymaster;