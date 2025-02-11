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

1. Изменить `bcrypt` на что-то другое, чтобы повысить производительность
2. Написать интеграционные тесты
3. Почистить код, переименовать некоторые методы/функции/переменные
