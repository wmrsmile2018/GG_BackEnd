red=`tput setaf 1`
green=`tput setaf 2`
yellow=`tput setaf 3`
purple=`tput setaf 5`
reset=`tput sgr0`

echo "${yellow}drop all tables and types${purple}"
PGGSSENCMODE=disable psql -h localhost -d ggaming -U wMrSmile -p 5432 -a -q -f /Users/wMrSmile/go/src/github.com/wmrsmile2018/GG/sql_commands/script_drop.sql
echo "${reset}"
