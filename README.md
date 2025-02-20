# Merch Shop (Avito Trainee Task)

[![Go Report Card](https://goreportcard.com/badge/github.com/0x0FACED/merch-shop?random=1)](https://goreportcard.com/report/github.com/0x0FACED/merch-shop)


Это мое решение тестового задания в Avito на позицию Backend Developer.

## Table of Contents

- [Merch Shop (Avito Trainee Task)](#merch-shop-avito-trainee-task)
  - [Table of Contents](#table-of-contents)
  - [Подход к решению задачи](#подход-к-решению-задачи)
  - [Архитектура сервера](#архитектура-сервера)
  - [Обработка ошибок](#обработка-ошибок)
  - [Проектирование базы данных](#проектирование-базы-данных)
  - [Установка и запуск](#установка-и-запуск)
    - [Запуск в docker](#запуск-в-docker)
    - [Запуск через go run/make](#запуск-через-go-runmake)
  - [Тестирование](#тестирование)
    - [Нагрузочное тестирование](#нагрузочное-тестирование)
      - [Load Test 1](#load-test-1)
      - [Load Test 2](#load-test-2)
      - [Load Test 3](#load-test-3)
      - [Load Test 4](#load-test-4)
      - [Load Test 5](#load-test-5)
      - [Load Test 6](#load-test-6)
    - [Профилирование во время нагрузочных тестов и после](#профилирование-во-время-нагрузочных-тестов-и-после)
    - [Unit-тесты](#unit-тесты)
    - [Интеграционное тестирование](#интеграционное-тестирование)
    - [Ручное тестирование](#ручное-тестирование)
      - [POST /api/auth](#post-apiauth)
      - [POST /api/info](#post-apiinfo)
      - [POST /api/sendCoin](#post-apisendcoin)
      - [GET /api/buy/:item](#get-apibuyitem)
  - [TODO](#todo)
  - [Вопросы и размышления](#вопросы-и-размышления)
    - [Вопросы к API](#вопросы-к-api)
    - [Конфигурация](#конфигурация)
    - [Хранение паролей](#хранение-паролей)
    - [Кэширование](#кэширование)
    - [Индексы](#индексы)
  - [Использованные технологии](#использованные-технологии)

## Подход к решению задачи

Увидев тестовое в первый раз, я удивился: "Какое-то оно очень простое и маленькое по сравнению с прошлым (Осень 2024)". Но потом увидел условия, а там тесты.

Решил не использовать кодогенерацию, так как API небольшой, все быстренько руками сам написал. 

Для написания API думал взять даже `net/http`, но потом подумал, что лучше `echo` использую. 

Для работы с базой данных использовал `pgx`, это я даже не обдумывал, а сразу решил.

## Архитектура сервера

Я решил, что архитектура должна быть 3-х уровневая: 

1. API
2. Service
3. Database

Запросы поступают в API, биндятся к структурам, отражающим входные данные, валидируются при помощи `go-playground/validator/v10` и далее, если все хорошо, данные запросов передаются в `service layer` в структурах, в названии которых в конце есть `Params`.

Стрктуры с окончанием `Params` несут в себе смысл параметров запросов, то есть там данные тел самих запросов. Это сделано для разделения уровней.

В `service layer` поступают уже валидированные данные. То есть, если отправить запрос на отправку монет и указать там `amount < 0`, то такой запрос не дойдет даже до уровня `service`. Здесь данные передаются в `database layer`. **Можно спросить: а зачем нужен `service layer`, если он как таковой функции не выполняет?** А нужен этот слой, потому что этой слой БЛ (бизнес логики). Если придется добавлять, например, работу с кэшированием через `redis` или запросы к внешним API, то их правильно будет расположить на этом слое. Тогда код `API layer` и `database layer` не будет изменен. Да и эти слои не должны этим заниматься.

В `database layer` выполняются запросы к базе `PostgreSQL`, результат возвращается в `service layer`, а оттуда в `API layer`.

## Обработка ошибок

Я этому моменту уделил довольно много времени. Решил остановиться на довольно простом варианте:

1. Есть ошибки уровня базы данных, сервисного уровня.
2. База возвращает обернутые ошибки `err` с заранее опредленным ошибками моими.
3. Сервис имеет аналогичный набор ошибок заранее определенных. Он маппит ошибки уровня базы в свои и отдают уже свои ошибки.
4. API маппит ошибки сервиса в `http` коды и отдает всегда `status code`.

Это довольно тривиальный подход, но весь маппинг зато вынесем в отдельные 2 функции, а в `API layer` и в `service layer` остается только вызвать функцию маппинга.

`Wrap` на уровне базы есть, чтобы в `service layer` можно было залоггировать ошибку полную, а в API отдать только самое важное.

## Проектирование базы данных

При проектировании я отталкивался от возможных сущностей и от сущностей, описанных в спецификации.

Я создал следующие таблицы:

1. `users`
2. `wallets`
3. `transactions`
4. `items`
5. `inventory`

`users` отвечает за хранение информации о пользователе.
`wallets` хранит кошельки пользователей и создается в момент создания пользователя автоматически
`transactions` хранит в себе транзакции между пользователями, но не хранит операции о покупках вещей.
`items` просто хранит все указанные в описании задания предметы и их стоимость.
`inventory` представляет из себя инвентарь пользователя, а именно предмет и количество этого предмета у конкретного пользователя по его `ID`.

Ниже приведена диаграмма полученной БД:

![DB Diagram](/images/db_diagram.png)

## Установка и запуск

Изначально необходимо склонировать репозиторий:

```sh
git clone https://github.com/0x0FACED/merch-shop.git
```


Что ж, есть несколько базовых способов запустить сервис.

1. **docker**
2. **go run** или **make**


### Запуск в docker

В корне проекта уже лежат файлы `Dockerfile` и `docker-compose.yml`. Осталось только по примеру из `.env.example` создать свой `.env` файлик и указать нужные параметры.

После этого можно запускать проект с помощью команды:

```sh
sudo docker compose up --build
```

### Запуск через go run/make

Если вдруг через Docker неудобно, то можно воспользоваться стандартным запуском через билд проекта.

Надо создать `.env` файлик по примеру `.env.example`, как и в случае с запуском через docker.

Если имеется **make**, то можно запустить через

```sh
make build-run
```

Эта команда забилдит executable файлик и запустит его.

## Тестирование

### Нагрузочное тестирование

Для проведения нагрузочных тестов использовалась утилита `k6`, являющаяся одним ищ лучших решений для такой задачи.

Установка до боли проста на моей системе:
```sh
yay -S k6
```

Все тесты должны писаться в `*.js` файлах, а вот это уже вызвало некоторое затруднение в виду моего практически полного незнания `js`.

#### Load Test 1

Для первого теста был написан файлик `load_test.js`. Все параметры в нем указаны, как и в выводе результатов ниже. 

**Тестовый файл**: `load_test.js`

*Тест запускается такой командой:*

```sh
k6 run load_test.js
```

*Результаты тестирования:*

```sh
         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: load_test.js
        output: -

     scenarios: (100.00%) 1 scenario, 500 max VUs, 1m30s max duration (incl. graceful stop):
              * default: 500 looping VUs for 1m0s (gracefulStop: 30s)


     ✓ Info success
     ✓ Buy success
     ✓ SendCoin success

     █ setup

       ✓ Auth success

   ✓ checks.........................: 100.00% 59839 out of 59839
     data_received..................: 11 MB   180 kB/s
     data_sent......................: 16 MB   258 kB/s
     http_req_blocked...............: avg=11.62µs  min=1.37µs   med=4.79µs   max=26.59ms p(90)=8.06µs   p(95)=9.69µs  
     http_req_connecting............: avg=5.62µs   min=0s       med=0s       max=26.54ms p(90)=0s       p(95)=0s      
   ✓ http_req_duration..............: avg=583.04µs min=115.38µs med=561.28µs max=76.59ms p(90)=900.69µs p(95)=1.08ms  
       { expected_response:true }...: avg=723.94µs min=338.41µs med=618.42µs max=76.59ms p(90)=972.35µs p(95)=1.16ms  
     http_req_failed................: 66.66%  39892 out of 59839
     http_req_receiving.............: avg=44.66µs  min=8.1µs    med=38.79µs  max=4.95ms  p(90)=69.04µs  p(95)=80.99µs 
     http_req_sending...............: avg=16.61µs  min=3.6µs    med=14.43µs  max=1.3ms   p(90)=25.18µs  p(95)=30.58µs 
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s      p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=521.76µs min=80.51µs  med=509.57µs max=76.38ms p(90)=818.91µs p(95)=982.43µs
     http_reqs......................: 59839   974.272818/s
     iteration_duration.............: avg=1.52s    min=1s       med=1.5s     max=2.49s   p(90)=1.5s     p(95)=1.67s   
     iterations.....................: 19946   324.752179/s
     vus............................: 216     min=216            max=500
     vus_max........................: 500     min=500            max=500


running (1m01.4s), 000/500 VUs, 19946 complete and 0 interrupted iterations
default ✓ [======================================] 500 VUs  1m0s
```

Стоит пояснить, почему `http_req_failed` **66%**, то есть всего **34%** тестов были успешны. Дело в том, что `http_req_failed` **считает за проваленный любой тест, который вернул код, отличающийся от** `2xx` или `3xx`. А, так как при создании пользователя у него имеется всего 1000 монет, то, например, отправлять 10000 раз по 1 монете и получать **200** код, **не получится**, потому что API отдает **400** код с ошибкой `"insuffisient funds"`. Поэтому было принято решение добавить в `thresholds` метрику `checks` и добавить в `check` **400** код как валидный, **потому что это ожидаемый ответ от сервера при отсутствии средств**.

А **34%** успешных тестов, потому что 1 тест - `/api/auth`, а остальные это `/api/info`, то есть получение информации, а в этих тестах все корректно отработало.

#### Load Test 2

В этом тесте я решил сделать иначе. Я тестирую только отправку монет между юзерами.

Создается 2 юзера, и они по очереди отправляют друг другу по 1 монете

**Тестовый файл**: `load_test_send_coin.js`

*Тест запускается такой командой:*

```sh
k6 run load_test_send_coin.js
```

*Результаты теста:*

```sh
         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: load_test_send_coin.js
        output: -

     scenarios: (100.00%) 1 scenario, 500 max VUs, 1m30s max duration (incl. graceful stop):
              * default: 500 looping VUs for 1m0s (gracefulStop: 30s)


     ✗ SendCoin success
      ↳  65% — ✓ 19628 / ✗ 10372

     █ setup

       ✓ Auth loadtest2user1 success
       ✓ Auth loadtest2user2 success

   ✗ checks.........................: 65.42% 19630 out of 30002
     data_received..................: 3.2 MB 52 kB/s
     data_sent......................: 9.8 MB 161 kB/s
     http_req_blocked...............: avg=17.84µs min=1.42µs   med=6.46µs  max=19.6ms  p(90)=10.71µs p(95)=13.11µs
     http_req_connecting............: avg=8.8µs   min=0s       med=0s      max=19.49ms p(90)=0s      p(95)=0s     
   ✓ http_req_duration..............: avg=2.2ms   min=766.19µs med=1.7ms   max=99.3ms  p(90)=2.71ms  p(95)=3.84ms 
       { expected_response:true }...: avg=1.85ms  min=1.04ms   med=1.63ms  max=79.88ms p(90)=2.39ms  p(95)=2.86ms 
   ✗ http_req_failed................: 34.57% 10372 out of 30002
     http_req_receiving.............: avg=54µs    min=8.9µs    med=45.84µs max=1.8ms   p(90)=78.41µs p(95)=98.53µs
     http_req_sending...............: avg=28.69µs min=5.8µs    med=24.97µs max=7ms     p(90)=41.38µs p(95)=49.68µs
     http_req_tls_handshaking.......: avg=0s      min=0s       med=0s      max=0s      p(90)=0s      p(95)=0s     
     http_req_waiting...............: avg=2.12ms  min=692.93µs med=1.62ms  max=99.19ms p(90)=2.6ms   p(95)=3.7ms  
     http_reqs......................: 30002  491.80274/s
     iteration_duration.............: avg=1s      min=1s       med=1s      max=1.49s   p(90)=1s      p(95)=1.01s  
     iterations.....................: 30000  491.769955/s
     vus............................: 66     min=66             max=500
     vus_max........................: 500    min=500            max=500


running (1m01.0s), 000/500 VUs, 30000 complete and 0 interrupted iterations
default ✓ [======================================] 500 VUs  1m0s
ERRO[0061] thresholds on metrics 'checks, http_req_failed' have been crossed
```

Можно заметить, что результаты гораздо хуже. За корректный ответ считался только код 200.

Это можно легко объяснить, так как отправка монет выполняется в рамках транзакции, поэтому присутствует блокировка данных. Из-за этого нельзя было получить доступ к ним, пока другая транзакция не завершится. Уровен изоляции: `Serializable`, что гарантирует отсутствие "грязного" чтения, неповторяемого чтения, фантомного чтения и аномалий сериализации. Это особенно важно, когда мы говорим про денежные переводы. 

Это означает, что многие попытки перевода были отклонены сервером в виду того, что была открыта какая-то другая транзакция.

**Самое важное, как я считаю:** изначально у одного юзера было 1000 монет, у второго 1000 монет, а по итогам тестирования у первого стало 1223, второго 777, **а значит никаких аномалий с деньгами не произошло**. Просто многие переводы были отклонены.

#### Load Test 3

Этот тест был направлен только на запрос `/api/info`.

**Тестовый файл**: `load_test_info.js`

Запускать командой:

```sh
k6 run load_test_info.js
```

```sh
         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: load_test_info.js
        output: -

     scenarios: (100.00%) 1 scenario, 500 max VUs, 1m30s max duration (incl. graceful stop):
              * default: 500 looping VUs for 1m0s (gracefulStop: 30s)


     ✓ Info success

     █ setup

       ✓ Auth success

   ✓ checks.........................: 100.00% 30001 out of 30001
     data_received..................: 5.9 MB  98 kB/s
     data_sent......................: 7.0 MB  115 kB/s
     http_req_blocked...............: avg=18.11µs  min=1.22µs   med=4.72µs   max=30.03ms p(90)=8.06µs  p(95)=10.05µs
     http_req_connecting............: avg=11.73µs  min=0s       med=0s       max=29.91ms p(90)=0s      p(95)=0s     
   ✓ http_req_duration..............: avg=852.37µs min=421.01µs med=718.72µs max=72.75ms p(90)=1.09ms  p(95)=1.3ms  
       { expected_response:true }...: avg=852.37µs min=421.01µs med=718.72µs max=72.75ms p(90)=1.09ms  p(95)=1.3ms  
   ✓ http_req_failed................: 0.00%   0 out of 30001
     http_req_receiving.............: avg=47.33µs  min=8.31µs   med=41.97µs  max=2.96ms  p(90)=69.02µs p(95)=81.68µs
     http_req_sending...............: avg=15.14µs  min=4.06µs   med=12.51µs  max=1.4ms   p(90)=22µs    p(95)=27.39µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s      p(90)=0s      p(95)=0s     
     http_req_waiting...............: avg=789.89µs min=392.15µs med=660.34µs max=71.8ms  p(90)=1.01ms  p(95)=1.21ms 
     http_reqs......................: 30001   493.537802/s
     iteration_duration.............: avg=1s       min=1s       med=1s       max=1.5s    p(90)=1s      p(95)=1s     
     iterations.....................: 30000   493.521352/s
     vus............................: 500     min=500            max=500
     vus_max........................: 500     min=500            max=500


running (1m00.8s), 000/500 VUs, 30000 complete and 0 interrupted iterations
default ✓ [======================================] 500 VUs  1m0s
```

Результаты показывают, что **все запросы оказались успешными, а среднее время задержки меньше 1мс**.


#### Load Test 4

Меня немного не устроил 3-й тест, поэтому я решил повысить количество виртуальных юзеров с 500 до 1000.

```sh
         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: load_test_info.js
        output: -

     scenarios: (100.00%) 1 scenario, 1000 max VUs, 1m30s max duration (incl. graceful stop):
              * default: 1000 looping VUs for 1m0s (gracefulStop: 30s)


     ✓ Info success

     █ setup

       ✓ Auth success

   ✓ checks.........................: 100.00% 59632 out of 59632
     data_received..................: 12 MB   193 kB/s
     data_sent......................: 14 MB   229 kB/s
     http_req_blocked...............: avg=24.19µs min=1.45µs   med=6.86µs   max=36.01ms p(90)=10.48µs  p(95)=12.9µs  
     http_req_connecting............: avg=14.98µs min=0s       med=0s       max=35.95ms p(90)=0s       p(95)=0s      
   ✓ http_req_duration..............: avg=1.05ms  min=439.46µs med=908.84µs max=56.79ms p(90)=1.37ms   p(95)=1.84ms  
       { expected_response:true }...: avg=1.05ms  min=439.46µs med=908.84µs max=56.79ms p(90)=1.37ms   p(95)=1.84ms  
   ✓ http_req_failed................: 0.00%   0 out of 59632
     http_req_receiving.............: avg=77.41µs min=12.41µs  med=65.75µs  max=7.21ms  p(90)=104.55µs p(95)=133.36µs
     http_req_sending...............: avg=22.13µs min=5.02µs   med=18.21µs  max=6.17ms  p(90)=27.7µs   p(95)=34.97µs 
     http_req_tls_handshaking.......: avg=0s      min=0s       med=0s       max=0s      p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=955.8µs min=387.44µs med=818.34µs max=56.38ms p(90)=1.24ms   p(95)=1.68ms  
     http_reqs......................: 59632   976.391119/s
     iteration_duration.............: avg=1.01s   min=1s       med=1s       max=1.99s   p(90)=1.01s    p(95)=1.01s   
     iterations.....................: 59631   976.374746/s
     vus............................: 175     min=175            max=1000
     vus_max........................: 1000    min=1000           max=1000
```

`http_req_duration` остался на том же примерно уровне.


#### Load Test 5

Этот тест проводится с заранее созданными 100к пользователями в базе данных.

**Как создать 100к юзеров в БД (Docker):**

1. Сначала надо запустить docker контейнеры:

```sh
sudo docker compose up -d
```

`-d` запускает контейнеры в фоновом режиме.

2. Копировать SQL скрипт для создания юзеров и кошельков в докер контейнер следующей командой:

```sh
sudo docker cp init_users_100k.sql db-container-name:/init_users_100k.sql
```

У меня контейнер называется `avito-shop-db`, поэтому на место `db-container-name` я пишу `avito-shop-db`.

3. Выполнить скрипт внутри контейнера. Скрипт создаст расширение `pgcrypto` для использования `bcrypt` и в цикле будет создавать 100к юзеров и 100к кошельков для них.

```sh
sudo docker exec -it db-container-name psql -U username -d db_name -f /init_users_100k.sql
```

Скрипт может выполняться продолжительное время.

Когда скрипт будет завершен, можно запускать нагрузочный тест на 100к юзеров. Скрипт читает файл `users.txt` и случайным образом берет из юзеров одного для выполнения запроса. Мы заранее создаем 100к этих юзеров, потому что создание юзера идет долго из-за вычисления хэша пароля.

**Тестовый файл**: `load_test_100k_users.js`

**Запуск теста**:

```sh
k6 run load_test_100k_users.js
```
**Результаты тестирования:**

```sh
         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: load_test_100k_users.js
        output: -

     scenarios: (100.00%) 1 scenario, 500 max VUs, 1m30s max duration (incl. graceful stop):
              * default: 500 looping VUs for 1m0s (gracefulStop: 30s)


     ✓ Auth success
     ✓ Info success

   ✓ checks.........................: 100.00% 58726 out of 58726
     data_received..................: 13 MB   208 kB/s
     data_sent......................: 13 MB   205 kB/s
     http_req_blocked...............: avg=13.88µs min=992ns    med=5.06µs  max=34.25ms p(90)=7.88µs  p(95)=9.36µs 
     http_req_connecting............: avg=7.58µs  min=0s       med=0s      max=34.19ms p(90)=0s      p(95)=0s     
   ✓ http_req_duration..............: avg=1.24ms  min=450.99µs med=1.34ms  max=13.76ms p(90)=1.82ms  p(95)=2.05ms 
       { expected_response:true }...: avg=1.24ms  min=450.99µs med=1.34ms  max=13.76ms p(90)=1.82ms  p(95)=2.05ms 
     http_req_failed................: 0.00%   0 out of 58726
     http_req_receiving.............: avg=48.47µs min=10.38µs  med=44.46µs max=3.62ms  p(90)=67.88µs p(95)=79.02µs
     http_req_sending...............: avg=18.45µs min=3.79µs   med=16.15µs max=4.5ms   p(90)=25.64µs p(95)=30.85µs
     http_req_tls_handshaking.......: avg=0s      min=0s       med=0s      max=0s      p(90)=0s      p(95)=0s     
     http_req_waiting...............: avg=1.17ms  min=409.27µs med=1.29ms  max=13.63ms p(90)=1.73ms  p(95)=1.97ms 
     http_reqs......................: 58726   962.659296/s
     iteration_duration.............: avg=1.03s   min=1s       med=1s      max=1.99s   p(90)=1.01s   p(95)=1.07s  
     iterations.....................: 29363   481.329648/s
     vus............................: 31      min=31             max=500
     vus_max........................: 500     min=500            max=500


running (1m01.0s), 000/500 VUs, 29363 complete and 0 interrupted iterations
default ✓ [======================================] 500 VUs  1m0s
```

Здесь мы только делали `auth` и запрашивали информацию о себе. **Средняя продолжительность запроса 1.2мс**, ни одной ошибки нет.

#### Load Test 6

Немного изменим скрипт для нагрузочного тестирования по сравнению с Load Test 5. А именно: добавим еще и отправку монеты и покупку предмета юзером помимо получения информации о себе.

**Результаты тестирования**:

```sh
         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: load_test_100k_users.js
        output: -

     scenarios: (100.00%) 1 scenario, 500 max VUs, 1m30s max duration (incl. graceful stop):
              * default: 500 looping VUs for 1m0s (gracefulStop: 30s)


     ✓ Auth success
     ✓ Info success
     ✓ Buy success
     ✗ SendCoin success
      ↳  99% — ✓ 15025 / ✗ 5

   ✓ checks.........................: 99.99% 60115 out of 60120
     data_received..................: 9.0 MB 146 kB/s
     data_sent......................: 15 MB  243 kB/s
     http_req_blocked...............: avg=11.2µs  min=1.33µs   med=5.36µs  max=23.05ms p(90)=8.21µs  p(95)=9.77µs 
     http_req_connecting............: avg=4.78µs  min=0s       med=0s      max=22.97ms p(90)=0s      p(95)=0s     
   ✓ http_req_duration..............: avg=1.66ms  min=571.15µs med=1.55ms  max=34.46ms p(90)=2.1ms   p(95)=2.4ms  
       { expected_response:true }...: avg=1.66ms  min=571.15µs med=1.55ms  max=34.46ms p(90)=2.1ms   p(95)=2.4ms  
     http_req_failed................: 0.00%  5 out of 60120
     http_req_receiving.............: avg=45.44µs min=8.2µs    med=41.97µs max=2.29ms  p(90)=63.73µs p(95)=74.14µs
     http_req_sending...............: avg=19.09µs min=4.49µs   med=17.21µs max=1.28ms  p(90)=27.39µs p(95)=32.5µs 
     http_req_tls_handshaking.......: avg=0s      min=0s       med=0s      max=0s      p(90)=0s      p(95)=0s     
     http_req_waiting...............: avg=1.59ms  min=537.04µs med=1.5ms   max=34.39ms p(90)=2.02ms  p(95)=2.32ms 
     http_reqs......................: 60120  975.996904/s
     iteration_duration.............: avg=2.02s   min=1s       med=1.99s   max=2.98s   p(90)=2.07s   p(95)=2.46s  
     iterations.....................: 15030  243.999226/s
     vus............................: 258    min=258            max=500
     vus_max........................: 500    min=500            max=500


running (1m01.6s), 000/500 VUs, 15030 complete and 0 interrupted iterations
default ✓ [======================================] 500 VUs  1m0s
```

Оценим результаты. Продолжительность запроса 1.66мс, что удовлетворяет условию задачи. Запросы на покупку все прошли, 100% успешных. Запросы на отправку монет прошли почти все, за исключением 5. Это связано с тем, что некоторые запросы обращались к одним и тем же данным. Например, одновременно так получалось, что рандом выбирал одного юзера для выполнения запросов. Или же одновременно `user2` выбирался, которому шла отправка монеты. Это небольшая погрешность. 

Считаю нагрузочные тесты пройдены успешно.

### Профилирование во время нагрузочных тестов и после

Во время проведения нагрузочных тестов результаты профилирования были следующими:

```sh
Types of profiles available:
Count	Profile
101	allocs
0	block
0	cmdline
514	goroutine
101	heap
0	mutex
0	profile
22	threadcreate
0	trace
```

А после профилирования такие:

```sh
Types of profiles available:
Count	Profile
129	allocs
0	block
0	cmdline
9	goroutine
129	heap
0	mutex
0	profile
21	threadcreate
0	trace
```

Утечек горутин нет.

### Unit-тесты

Задачей было покрыть основные бизнес сценарии, то есть покрыть тестами основную бизнес логику. Соответственно, тесты писались для `service layer`, потому что это там и расположены все бизнес сценарии.

*Для проверки процента покрытия тестами выполним поманду:*
```sh
go test ./... -cover
```

**Покрытие тестами:**

```sh
?       github.com/0x0FACED/merch-shop/internal/model   [no test files]
        github.com/0x0FACED/merch-shop/internal/service/mocks           coverage: 0.0% of statements
        github.com/0x0FACED/merch-shop/internal/server/validator                coverage: 0.0% of statements
        github.com/0x0FACED/merch-shop/internal/server          coverage: 0.0% of statements
        github.com/0x0FACED/merch-shop/internal/database                coverage: 0.0% of statements
        github.com/0x0FACED/merch-shop/config           coverage: 0.0% of statements
        github.com/0x0FACED/merch-shop/cmd/app          coverage: 0.0% of statements
        github.com/0x0FACED/merch-shop/pkg/logger               coverage: 0.0% of statements
        github.com/0x0FACED/merch-shop/internal/server/handler          coverage: 0.0% of statements
ok      github.com/0x0FACED/merch-shop/internal/service (cached)        coverage: 86.2% of statements
ok      github.com/0x0FACED/merch-shop/tests/e2e        (cached)        coverage: 100.0% of statements
```

Нас интересует строка `ok      github.com/0x0FACED/merch-shop/internal/service (cached)        coverage: 86.2% of statements`.

В ней сказано, что **покрыто 86.2% сценариев тестами**, что удовлетворяет условию в минимум 40%.

Таким образом, основные бизнес сценарии были покрыты тестами.

Для мокирования базы данных использовалась библиотека `stretchr/testify/mock`. Было довольно просто даже без всяких `mockgen` написать руками `mock repo`. Если бы репозиторий был большой по количеству методов, то лучше было бы использовать кодогенерацию.

*Запустить тесты можно с помощью `make` команды:*

```sh
make run-tests
```

**Вот результат ее выполнения:**

```sh
go test -v ./internal/service > ./tests/service_tests.log 2>&1
go test -v ./tests/e2e > ./tests/e2e_tests.log 2>&1
```

Результаты записаны в файлы, все тесты успешны, ошибок не было.

### Интеграционное тестирование

Было решено написать интеграционные тесты для обработчиков API. Для этого была создана отдельная директория `./tests/e2e`. Это, конечно, не совсем `e2e` тесты, но директорию я назвал так, потому что `integration` довольно длинное слово. Но это именно интеграционные тесты.

Для проведения этих тестов необходимо создать тестовую базу данных. Я назвал ее `merch_shop_test`. Конечно, можно проводить их и на основной базе данных, но это как-то неправильно мягко говоря.

Для интеграционных тестов был написан отдельный `.env.test` файл, который копирует содержание `.env.example`, но базу надо указывать тестовую.

Миграции применить к этой базе надо заранее. У меня это выполняется в команде `make run-tests` автоматически.

Для загрузки тестового конфига была написана функция `MustLoadTestConfig()`, которая паникует, если не удалось загрузить конфиг. Путь к файлу необходимо прописать самостоятельно, так как просто передав `.env.test` файл не будет найден. Но можно положит этот `.env.test` рядом с функцией, то есть в `./config/.env.test`. Тогда можно будет не указывать путь.

### Ручное тестирование

Для ручного тестирования использовалась программа `Postman`.

#### POST /api/auth

Попробуем аутентифицироваться в первый раз в системе:

![Auth Test3 User](/images/image.png)

Получаем в ответ токен, который будем использовать в дальнейших запросах. 

В базе уже имеются юзеры `test1`, `test2`.

#### POST /api/info

Получим информацию о себе.

![alt text](/images/image-1.png)

Как мы видим, у `test3` имеется 1000 монет, пустой инвентарь и не операций передачи монет и получения монет. Это логично, так как юзер только что был создан.

#### POST /api/sendCoin

Отправим 50 монет другому юзеру, а именно юзеру `test2`.

![alt text](/images/image-2.png)

И после этого запросим информацию о себе.

![alt text](/images/image-3.png)

Как мы видим, количество монет уменьшилось, а так же получилось исходящая транзакция.

#### GET /api/buy/:item

Купим несколько предметов, например, `pen`.

![alt text](/images/image-4.png)

Мы купили 4 раза ручку, теперь запросим информацию о себе.

![alt text](/images/image-5.png)

Окей, теперь у нас в инвентаре имеется еще 4 ручки, денег стало, соответственно на 40 монет меньше (10 монет одна ручка).

Купим еще один предмет: `hoody`.

![alt text](/images/image-6.png)

Теперь у нас есть еще и худак! И это прекрасно!

Вот было бы хорошо, если бы кто-то нам монеток отправил...

![alt text](/images/image-7.png)

О как прекрасно! Нам еще и `test1` отправил 500 монет и 250 монет!

Тогда можем купить еще несколько предметов, например, `powerbank` и `pink-hoody`. Будем самыми модными в офисе!

![alt text](/images/image-8.png)

Замечательно! Теперь у нас есть еще и повербанк + самый модный худак.

## TODO

- [x] Спроектировать архитектуру сервиса
- [x] Написать методы для работы с базой
- [x] Написать методы сервисного уровня
- [x] Написать API обработчики
- [ ] Изменить `bcrypt` на что-то другое, чтобы повысить производительность (опционально)
- [x] Написать интеграционные тесты
- [x] Написать unit-тесты
- [x] Почистить код, переименовать некоторые методы/функции/переменные
- [x] Написать `Dockerfile` и `docker-compose.yml` файлы
- [x] Сделать `doc` комментарии у пакетов (опционально)
- [x] Перепроектировать маппинг ошибок от базы к API (опционально)
- [x] Написать полный `README.md`
- [x] Добавить еще линтеров (опционально)
- [x] Провести нагрузочное тестирование, добавить в `README.md`
- [x] Добавить профилирование
- [x] Добавить теги `validate`
- [x] Убрать базу из сервера
- [x] Прокидывать секрет для JWT где-то извне
- [x] Глобальный рефакторинг
- [x] Добавить нормальные логи
- [ ] Добавить кэширование (`redis`) (опционально)
- [ ] Добавить Github Actions (опционально)

## Вопросы и размышления

Здесь находятся все мои возникшие вопросы в ходе выполнения задания, а так же мои предложения.

### Вопросы к API

Самый главный вопрос, который появился моментально после открытия спецификации `open api`: **Почему эндпоинт `/api/buy/:item` `GET`, а не `POST`?**. Ведь **это НЕ идемпотентная операция**. Мы не можем делать такие запросы бесконечно и ожидать всегда одинаковый результат, так как после каждого запроса мы изменяем состояние сервера, а именно базы данных: мы уменьшаем наш баланс и добавляем себе в инвентарь новый предмет. А когда наш баланс будет слишком низок дял покупки предмета, то мы получми отказ от сервера. 

GET запросы идемпотентны, они используются же только для получения данных (должны использоваться), **GET не должен менять состояние системы**.

Я бы сделал этот запрос как `POST /api/items/:item/buy`. Так мы заранее определяем, что возможно горизонтальное расширение для предметов: подарить предмет, продать предмет и тд. Например, `POST /api/items/:item/gift` или `POST /api/items/:item/sell`.

Аналогично с отправкой монет я бы сделал `POST /api/wallets/:wallet_id/send`. Тело запроса останется таким же, а `wallet_id` - это id нашего кошелька. Тогда горизонтально можно расшириться до следующего:

`POST /api/wallets/:wallet_id/deposit` - пополнить
`POST /api/wallets/:wallet_id/withdraw` - вывести
`POST /api/wallets/:wallet_id/balance` - посмотреть баланс

### Конфигурация

Я не совсем понял из описания задания, как требуется оформить конфигурацию сервера: при запуске докера ее прокидывать? Или из `.env` файла грузить в переменные окружения докера?

В итоге я сделал `.env` файл, где есть все необходимое. Таким образом, я конфигурирую приложение полностью, в том числе прокидываю `JWT_SECRET_KEY`.

### Хранение паролей

В задании не было ничего сказано, что нельзя хранить пароли в `plain text`, но я не могу хранить их так, поэтому добавить получение хэша пароля и хранение именно хэша с солью. Алгоритм `bcrypt` - один из самых простых и базовых. В продакшине можно на что-то получше заменить.

### Кэширование

У меня еще есть мысли добавить кэширование для часто запрашиваемых данных. Если добавлять, то я бы создал пакет `./internal/cache`, где поместил бы клиент `Redis`. Указатель на созданного клиента хранился бы в `service layer`, тогда, например, при запросе информации о себе запрос из API попадал бы в `service`, а здесь была проверка нахождения нужных данных в кэше. Если есть - достаем и не идем в базу. Если нет - идем в базу, достаем из базы, сохраняем в кэш и отдаем в API. Сделал бы какой-нибудь `TTL` на уровне хотя бы минут 30, наверное. Хотя для разных ключей можно разный `TTL` выставить в целом. 

### Индексы

Я проверил использование индексов через `EXPLAIN ANALYZE`. 

Результаты:

```sh
merch_shop=# EXPLAIN ANALYZE SELECT * FROM shop.users WHERE username = 'test_user';
                                                         QUERY PLAN                
                                         
-----------------------------------------------------------------------------------
-----------------------------------------
 Index Scan using idx_users_username on users  (cost=0.14..8.16 rows=1 width=560) (
actual time=0.007..0.007 rows=0 loops=1)
   Index Cond: ((username)::text = 'test_user'::text)
 Planning Time: 0.580 ms
 Execution Time: 0.036 ms
(4 rows)

merch_shop=# EXPLAIN ANALYZE 
SELECT * FROM shop.transactions 
WHERE from_user_id = 1 OR to_user_id = 1;
                                                              QUERY PLAN           
                                                    
-----------------------------------------------------------------------------------
----------------------------------------------------
 Bitmap Heap Scan on transactions  (cost=8.43..19.06 rows=16 width=24) (actual time
=0.007..0.008 rows=0 loops=1)
   Recheck Cond: ((from_user_id = 1) OR (to_user_id = 1))
   ->  BitmapOr  (cost=8.43..8.43 rows=16 width=0) (actual time=0.003..0.004 rows=0
 loops=1)
         ->  Bitmap Index Scan on idx_transactions_from_to  (cost=0.00..4.21 rows=8
 width=0) (actual time=0.002..0.002 rows=0 loops=1)
               Index Cond: (from_user_id = 1)
         ->  Bitmap Index Scan on idx_transactions_to_user  (cost=0.00..4.21 rows=8
 width=0) (actual time=0.001..0.001 rows=0 loops=1)
               Index Cond: (to_user_id = 1)
 Planning Time: 0.442 ms
 Execution Time: 0.053 ms
(9 rows)

merch_shop=# EXPLAIN ANALYZE 
SELECT * FROM shop.wallets 
WHERE balance > 500;
                                                          QUERY PLAN               
                                           
-----------------------------------------------------------------------------------
-------------------------------------------
 Bitmap Heap Scan on wallets  (cost=9.99..29.40 rows=753 width=8) (actual time=0.00
4..0.004 rows=0 loops=1)
   Recheck Cond: (balance > 500)
   ->  Bitmap Index Scan on idx_wallets_balance  (cost=0.00..9.80 rows=753 width=0)
 (actual time=0.001..0.001 rows=0 loops=1)
         Index Cond: (balance > 500)
 Planning Time: 0.162 ms
 Execution Time: 0.014 ms
(6 rows)

merch_shop=# EXPLAIN ANALYZE 
SELECT * FROM shop.inventory 
WHERE item_id = 6;
                                                         QUERY PLAN                
                                         
-----------------------------------------------------------------------------------
-----------------------------------------
 Bitmap Heap Scan on inventory  (cost=4.23..14.79 rows=10 width=12) (actual time=0.
004..0.004 rows=0 loops=1)
   Recheck Cond: (item_id = 6)
   ->  Bitmap Index Scan on idx_inventory_item  (cost=0.00..4.23 rows=10 width=0) (
actual time=0.001..0.001 rows=0 loops=1)
         Index Cond: (item_id = 6)
 Planning Time: 0.421 ms
 Execution Time: 0.014 ms
(6 rows)

merch_shop=# EXPLAIN ANALYZE 
UPDATE shop.inventory 
SET quantity = quantity + 1 
WHERE user_id = 1 AND item_id = 3;
                                                                    QUERY PLAN     
                                                                
-----------------------------------------------------------------------------------
----------------------------------------------------------------
 Update on inventory  (cost=0.15..8.17 rows=0 width=0) (actual time=0.005..0.006 ro
ws=0 loops=1)
   ->  Index Scan using idx_inventory_user_item_quantity on inventory  (cost=0.15..
8.17 rows=1 width=10) (actual time=0.003..0.003 rows=0 loops=1)
         Index Cond: ((user_id = 1) AND (item_id = 3))
 Planning Time: 0.108 ms
 Execution Time: 0.111 ms
(5 rows)
```

Запросы были оптимизированы для работы с индексами и Bitmap-сканированием, что позволяет им работать с большими объемами данных эффективно. Время выполнения большинства запросов составляет миллисекунды.

## Использованные технологии

1. [env](https://github.com/caarlos0/env)
2. [validator](https://github.com/go-playground/validator)
3. [jwt](https://github.com/golang-jwt/jwt)
4. [pgx](https://github.com/jackc/pgx)
5. [godotenv](https://github.com/joho/godotenv)
6. [echo](https://github.com/labstack/echo)
7. [testify](https://github.com/stretchr/testify)
8. [zap](https://go.uber.org/zap)