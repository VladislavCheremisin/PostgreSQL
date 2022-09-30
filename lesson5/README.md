В данном проекте мы можем использовать два варианта инициализации и наполнения данными БД.

1 Вариант.
Использовать специальные файлы при запуске контейнера
0001_init.sql - инициализируем БД gopher_library
0002_init.sql - Добавляем таблицы
0003_init.sql - Наполняем таблицы

2 Вариант.
В этом варианте мы будем использовать утилиту github.com/golang-migrate/migrate, она позволит нам производить Миграции БД.
Для запуска с использованием этого варианта, необходимо очистить файлы 0002_init.sql и 0003_init.sql.

Далее после запуска контейнера произвести команду
migrate -database "postgresql://postgres:P@ssw0rd@localhost:5432/gopher_library?sslmode=disable" -path migrations up 1  
Произойдет добавление таблиц в нашу БД.  
------------------------------------------------------------------------------------------
In this project, we can use two options for initializing and filling the database with data.

1 option.
Use special files when starting a container
0001_init.sql - initialize the database gopher_library
0002_init.sql - Add tables
0003_init.sql - Populate tables

Option 2.
In this option, we will use the github.com/golang-migrate/migrate utility, it will allow us to perform DB Migrations.

To start using this option, you need to clear the 0002_init.sql and 0003_init.sql files.
Next, after starting the container, issue the command
migrate -database "postgresql://postgres:P@ssw0rd@localhost:5432/gopher_library?sslmode=disable" -path migrations up 1
Tables will be added to our database.