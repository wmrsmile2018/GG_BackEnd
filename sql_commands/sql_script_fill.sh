red=`tput setaf 1`
green=`tput setaf 2`
yellow=`tput setaf 3`
purple=`tput setaf 5`
reset=`tput sgr0`

echo "${yellow}insert users${red}"
PGGSSENCMODE=disable psql -h localhost -d ggaming -U wMrSmile -p 5432 -a -q -f /Users/wMrSmile/go/src/github.com/wmrsmile2018/GG/sql_commands/tests/script_tests_insert_users.sql
echo "${yellow}insert balance${red}"
PGGSSENCMODE=disable psql -h localhost -d ggaming -U wMrSmile -p 5432 -a -q -f /Users/wMrSmile/go/src/github.com/wmrsmile2018/GG/sql_commands/tests/script_tests_insert_balance.sql
echo "${yellow}insert black_list${red}"
PGGSSENCMODE=disable psql -h localhost -d ggaming -U wMrSmile -p 5432 -a -q -f /Users/wMrSmile/go/src/github.com/wmrsmile2018/GG/sql_commands/tests/script_tests_insert_black_list.sql
echo "${yellow}insert chats${red}"
PGGSSENCMODE=disable psql -h localhost -d ggaming -U wMrSmile -p 5432 -a -q -f /Users/wMrSmile/go/src/github.com/wmrsmile2018/GG/sql_commands/tests/script_tests_insert_chats.sql
echo "${yellow}insert finance${red}"
PGGSSENCMODE=disable psql -h localhost -d ggaming -U wMrSmile -p 5432 -a -q -f /Users/wMrSmile/go/src/github.com/wmrsmile2018/GG/sql_commands/tests/script_tests_insert_finance.sql
echo "${yellow}insert friends${red}"
PGGSSENCMODE=disable psql -h localhost -d ggaming -U wMrSmile -p 5432 -a -q -f /Users/wMrSmile/go/src/github.com/wmrsmile2018/GG/sql_commands/tests/script_tests_insert_friends.sql
echo "${yellow}insert messages${red}"
PGGSSENCMODE=disable psql -h localhost -d ggaming -U wMrSmile -p 5432 -a -q -f /Users/wMrSmile/go/src/github.com/wmrsmile2018/GG/sql_commands/tests/script_tests_insert_messages.sql
echo "${yellow}insert commands${red}"
PGGSSENCMODE=disable psql -h localhost -d ggaming -U wMrSmile -p 5432 -a -q -f /Users/wMrSmile/go/src/github.com/wmrsmile2018/GG/sql_commands/tests/script_tests_insert_commands.sql
echo "${yellow}insert command_structure${red}"
PGGSSENCMODE=disable psql -h localhost -d ggaming -U wMrSmile -p 5432 -a -q -f /Users/wMrSmile/go/src/github.com/wmrsmile2018/GG/sql_commands/tests/script_tests_insert_command_structure.sql
echo "${yellow}insert histories${red}"
PGGSSENCMODE=disable psql -h localhost -d ggaming -U wMrSmile -p 5432 -a -q -f /Users/wMrSmile/go/src/github.com/wmrsmile2018/GG/sql_commands/tests/script_tests_insert_histories.sql
echo "${yellow}insert history_structure${red}"
PGGSSENCMODE=disable psql -h localhost -d ggaming -U wMrSmile -p 5432 -a -q -f /Users/wMrSmile/go/src/github.com/wmrsmile2018/GG/sql_commands/tests/script_tests_insert_history_structure.sql
echo "${reset}"
