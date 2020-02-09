\! tput setaf 1; echo "___________________________drop references____________________________";
\! tput setaf 2;
alter table friends drop constraint users_friends;
alter table black_list drop constraint users_black_list;
alter table chats drop constraint users_chats;
alter table balance drop constraint users_balance;
alter table finance drop constraint users_finance;
alter table commands drop constraint users_commands;
alter table history_structure drop constraint users_history_structure;
alter table messages drop constraint chats_messages;
alter table history_structure drop constraint histories_history_structure;


\! tput setaf 1; echo "_______________delete all rows from all tables _________________";
\! tput setaf 2;
delete from friends;
delete from black_list;
delete from balance;
delete from finance;
delete from messages;
delete from chats;
delete from command_structure;
delete from commands;
delete from history_structure;
delete from histories;
delete from users;


\! tput setaf 1; echo "___________________________drop all tables ____________________________";
\! tput setaf 2;
drop table black_list;
drop table friends;
drop table messages;
drop table chats;
drop table balance;
drop table finance;
drop table commands;
drop table command_structure;
drop table history_structure;
drop table histories;
drop table users;

\! tput setaf 1; echo "_________________________drop all types______________________________";
\! tput setaf 2;


drop type sex;
drop type type_chat;
drop type type_user;



/*  */
