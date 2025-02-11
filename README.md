# Merch Shop (Avito Trainee Task)

[![Go Report Card](https://goreportcard.com/badge/github.com/0x0FACED/merch-shop)](https://goreportcard.com/report/github.com/0x0FACED/merch-shop)


Это мое решение тестового задания в Avito на позицию Backend Developer.

## Table of Contents

- [Merch Shop (Avito Trainee Task)](#merch-shop-avito-trainee-task)
	- [Table of Contents](#table-of-contents)
	- [Подход к решению задачи](#подход-к-решению-задачи)
	- [TODO](#todo)

## Подход к решению задачи

Увидев тестовое в первый раз, я удивился: "Какое-то оно очень простое и маленькое по сравнению с прошлым (Осень 2024)". Но потом увидел условия, а там тесты.

Решил не использовать кодогенерацию, так как API небольшой, все быстренько руками сам написал. 

Для написания API думал взять даже `net/http`, но потом подумал, что лучше `echo` использую. 

Для работы с базой данных использовал `pgx`, это я даже не обдумывал, а сразу решил.

## TODO

- [x] Спроектировать архитектуру сервиса
- [x] Написать методы для работы с базой
- [x] Написать методы сервисного уровня
- [x] Написать API обработчики
- [ ] Изменить `bcrypt` на что-то другое, чтобы повысить производительность (опционально)
- [x] Написать интеграционные тесты
- [x] Написать unit-тесты
- [ ] Почистить код, переименовать некоторые методы/функции/переменные
- [ ] Написать `Dockerfile` и `docker-compose.yml` файлы
- [ ] Сделать `doc` комментарии у пакетов (опционально)
- [ ] Перепроектировать маппинг ошибок от базы к API (опционально)
- [ ] Написать полный `README.md`
- [ ] Добавить еще линтеров (опционально)
- [ ] Провести нагрузочное тестирование, добавить в `README.md`
- [ ] Добавить профилирование
