\! tput setaf 1; "___________________________create table users ____________________________";
\! tput setaf 2;


create table users (
  number serial not null,
  id_user varchar primary key,
  email varchar unique,
  type_user type_user,
  /* mob_number varchar, */
  login varchar,
  birthday timestamp,
  sex sex,
--   first_profit timestamp,
  encrypted_password varchar,
  /* status varchar, */
  /* avatar varchar, */
  creation_date timestamp,
  /* url_vk varchar, */
  /* url_facebook varchar, */
  /* url_twitter varchar, */
  /* url_instagram varchar, */
  /* url_twitch varchar, */
   url_youtube varchar
  );


-- \! tput setaf 3; "____________________________________________________________________________________";
-- \! tput setaf 1; "___________________________create table black_list ____________________________";
-- \! tput setaf 2;


create table black_list (
  number serial not null,
  id_user varchar,
  id_black_user varchar
  );


\! tput setaf 3; "____________________________________________________________________________________";
\! tput setaf 1; "___________________________create table friends ____________________________";
\! tput setaf 2;


create table friends (
  number serial not null,
  id_user varchar,
  id_friend varchar
  );


\! tput setaf 3; "____________________________________________________________________________________";
\! tput setaf 1; "___________________________create table messages ____________________________";
\! tput setaf 2;


create table messages (
  number serial not null,
  id_message varchar primary key,
  id_user varchar,
  id_chat varchar,
  text_mes text,
  creation_time timestamp,
  type_chat type_chat
  );


\! tput setaf 3; "____________________________________________________________________________________";
\! tput setaf 1; "___________________________create table chats ____________________________";
\! tput setaf 2;

create table users_chats (
    id_chat varchar,
    id_user varchar
);


create table chats (
  number serial not null,
  id_chat varchar primary key,
  id_user varchar,
  type_chat type_chat
);


\! tput setaf 3; "____________________________________________________________________________________";
\! tput setaf 1; "___________________________create table balance ____________________________";
\! tput setaf 2;


create table balance (
  id_user varchar primary key,
  BIC varchar,
  balance integer
  );


\! tput setaf 3; "____________________________________________________________________________________";
\! tput setaf 1; "___________________________create table finance ____________________________";
\! tput setaf 2;


create table finance (
  id_user varchar primary key,
  array_profits integer[],
  array_times timestamp[]
  );


\! tput setaf 3; "____________________________________________________________________________________";
\! tput setaf 1; "___________________________create table commands ____________________________";
\! tput setaf 2;


create table commands (
  number serial not null,
  id_command varchar primary key,
  id_user varchar unique not null
  );


\! tput setaf 3; "____________________________________________________________________________________";
\! tput setaf 1; "___________________________create table command_structure ____________________________";
\! tput setaf 2;


create table command_structure (
  id_command varchar,
  id_user varchar
  );


\! tput setaf 3; "____________________________________________________________________________________";
\! tput setaf 1; "___________________________create table histories ____________________________";
\! tput setaf 2;


create table histories (
  number serial not null,
  id_history varchar primary key,
  type_mode type_mode
  );


\! tput setaf 3; "____________________________________________________________________________________";
\! tput setaf 1; "___________________________create table history_structure ____________________________";
\! tput setaf 2;


create table history_structure (
  number serial not null,
  type_mode type_mode,
  id_history varchar,
  id_user varchar,
  id_command varchar,
  profit integer,
  is_winner bool,
  kills integer,
  deaths integer,
  score integer,
  sending_time timestamp,
  g_date date,
  url_video varchar
  );


\! tput setaf 3; "____________________________________________________________________________________";
