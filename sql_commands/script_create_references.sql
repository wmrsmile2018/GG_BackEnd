/* constraint parent_table: users */
alter table friends add constraint users_friends foreign key (id_user) references users (id_user);
alter table black_list add constraint users_black_list foreign key (id_user) references users (id_user);
alter table balance add constraint users_balance foreign key (id_user) references users (id_user);
alter table finance add constraint users_finance foreign key (id_user) references users (id_user);
alter table commands add constraint users_commands foreign key (id_user) references users (id_user);
alter table history_structure add constraint users_history_structure foreign key (id_user) references users (id_user);
alter table command_structure add constraint users_command_structure foreign key (id_user) references users (id_user);
alter table users_chats add constraint users_users_chats foreign key (id_user) references users (id_user);

/* constraint parent_table: chats */
alter table messages add constraint chats_messages foreign key (id_chat) references chats (id_chat);

/* constraint parent_table: histories */
alter table history_structure add constraint histories_history_structure foreign key (id_history) references histories (id_history);

/* constraint parent_table: commands */
alter table command_structure add constraint commands_command_structure foreign key (id_command) references commands (id_command);
