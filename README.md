
# Как запустить

## То с чем точно работает
go 1.17.3
mongodb 4.4.4
Linux Fedora 33

## Mongo
1. Добавить соответствующий репозиторий. Например не Федоре, это:
```
cd /etc/yum.repos.d
touch mongodb-org-4.4.repo
```
2. В файл `mongodb-org-4.4.repo` записать: 
```
[mongodb-org-4.4]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/8/mongodb-org/4.4/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-4.4.asc
```
`sudo dnf install mongodb-org`

3. Запустить сервер БД
`sudo mongod`

4. Проверить, что всё работает
`mongosh` -- Должна открыться оболочка командной строки

5. Альтернативно можно зарегистрироваться на сайте mongodb.com (нужен ВПН) и воспользоваться пробной версией их сервиса Atlas. Тогда перед использованием программы нужно будет подправить конфиг.

## Go
Смотри https://go.dev/doc/install

## Запуск
Клонируем репозиторий и переходим в папку с кодом.
```
git clone https://github.com/telephrag/bscrap
cd bscrap
``` 

Также можно просто скачать этот репозиторий, на этой странице.

Если вы создали пробную БД с помощью Atlas, то нужно в личном кабинете скопировать URI и заменить им значение константы `DBUri` в файле `bscrap/config/db.go`.

Скачиваем зависимости и запускаем
```
go mod vendor
go run main.go
``` 

В браузере переходим по ссылке: http://localhost:8080/?symbolA=BTCUSDT&symbolB=ZECUSDT&interval=1w&limit=52&startTime=1621728000000

В окне браузера появится результат обработки данных с API вместе с метаданными. Запоминаем значения полей `_id`, `_raw_data_id`.

Переходим по ссылке: http://localhost:8080/retrieve?processed={_id}&raw={_raw_data_id}, подставив соотвествующие значения.

В окне появится результат обработки данных с API вместе с данными, которые обрабатывались в формате удобном для БД из-за лимита по уровням вложенности у Монго. Данные будут взяты из БД. Данные удалятся СУБД автоматически по истечениею 1-2 минут (интервал обновления у Монго -- одна минута). 

Нажимаем `Ctrl+C`, чтобы завершить работу программы или отправляем дургой запрос.
