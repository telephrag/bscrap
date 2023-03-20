
# Что это
Это программа, загружающая исторические данные о ценах двух криптовалютных пар в заданном промежутке времени с заданными интервалом между записями и вычисляющая для них статистические данные такие как:
- мат.ожидание
- дисперсию
- корреляцию
- ковариацию

Исходные данные использованные для вычислений записываются в базу данных для возможного повторного использования в вычислених.Результаты вычислений также записываются в базу данных.

# Как запустить

## Зависимости, которые нужно установить самому
1. Docker 20.10.17, build 100c701 (скорее всего будет работать и с более поздними версиями версиями)
2. Docker Compose 2.16.0
3. Дистрибутив Linux, поддерживающий вышеназванные зависимости

## Зависимости, который Докер скачает сам
1. Golang 1.18.4
2. MongoDB 6.0.4
3. См. `go.mod`

## Запуск
Клонируем репозиторий и переходим в папку с кодом:
```
git clone https://github.com/telephrag/bscrap
cd bscrap
``` 
Возможно после установки Докера потребуется вручную запустить его сервис. На дистрибутивах с `systemd`
```
sudo systemctl start docker.service
```
Собираем и запускаем приложение:
```
sudo docker compose build && sudo docker compose up
```
Для повторного запуска в дальнейшем можно использовать просто `sudo docker compose up`.

# Использование
Запрос к сервису состоит из 6 параметров:

- `symbolA`, `symbolB` -- наименования криптовалютных пар. 
Список всех возможных пар можно посмотреть здесь: https://github.com/binance/binance-public-data/blob/master/data/symbols.txt

- `interval` -- интервал между замерами цен. 
Список всех возможных интервалов можно посмотреть здесь: https://github.com/binance/binance-public-data/#intervals 

- `limit` -- максимальное количество записей, которые можно загрузить

- `startTime`, `endTime` -- границы временного промежутка из которого можно брать данные. Для каждого интервала Бинанс хранит записи, сделанные в начале интервала. Например, если пользователь запрашивает записи с интервалом в неделю, но указывает начальное время, соответствующее, скажем, четвергу, то он получит набор записей первая из которого будет с делана в понедельник следующей недели.

`symbolA`, `symbolB`, `interval` -- обязательные параметры. 

По умолчанию:
- `startTime` = 0
- `endTime` = "момент запроса"
- `limit` = 500 

# Пример использования

Запустив сервис, в браузере переходим по ссылке: http://localhost:8080/?symbolA=BTCUSDT&symbolB=ZECUSDT&interval=1w&limit=52&startTime=1621728000000

Получив ответ, обновляем страницу и видим, что значение поля `fromDB` для обоих пар изменилось на `true`, а значения `_raw_data_?_id` остались неизменными. Это означает, что для вычислений были использованы данные уже имевшиеся в БД.

Запоминаем значения полей `_id`, `_raw_data_a_id`, `_raw_data_b_id`.

Переходим по ссылке: http://localhost:8080/retrieve?processed={_id}&rawA={_raw_data_a_id}&rawB={_raw_data_b_id}, подставив соотвествующие значения.

В окне появится результат обработки данных с API вместе с данными, которые обрабатывались в формате удобном для БД из-за лимита по уровням вложенности у Монго. Данные будут взяты из БД. Удалятся СУБД автоматически по истечениею 1-2 минут. 

Нажмите комбинацию `Ctrl+C` в окне терминала с приложением, чтобы завершить работу программы.
