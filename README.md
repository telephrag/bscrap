
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

5. Альтернативно можно зарегистрироваться на сайте mongodb.com (нужен ВПН) и воспользоваться пробной версией
их сервиса Atlas. Тогда перед использованием программы нужно будет подправить конфиг.

## Go
Смотри https://go.dev/doc/install

## Запуск
1. `git clone https://github.com/telephrag/bscrap` -- или просто скачать этот репозиторий, на этой странице
2. `cd bscrap` -- переход в директорию исходного кода
Если вы создали пробную БД с помощью Atlas, то нужно в личном кабинете скопировать URI и заменить им значение
константы `DBUri` в файле `bscrap/config/db.go`.
1. `go mod vendor` -- скачивание зависимостей 
2. `go run main.go` или `go build main.go`; `./main`. 
Через мгновение после запуска программа временно остановится.
3. `mongosh` -- открытие оболочки Монго
В оболоче Монго:
`use bscrapdb` -- переключение на используемую программой БД
`db.bscrap.find()` -- вывод всех документов коллекции, куда программа записывает результаты 
4. В окне терминала, где была запущена программа:
`Ctrl+C` -- продолжение работы программы, произойдёт удаление только что добавленного документа, в этом можно убедиться повторно прописав предыдущую команду в оболочке Монго