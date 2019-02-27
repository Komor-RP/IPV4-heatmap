# IPV4-heatmap

## Development Notes
Running sql file on the command line
```shell
psql -U postgres -d postgres -a -f ./loadData.sql
```
Running sql file on heroku
```shell
heroku pg:psql -c "command" --app "name-app"
```