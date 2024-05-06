# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

Сборка образа Docker:
`docker build --tag finalproject .`

Для запуска контейнера необходимо открыть терминал в текущей директории и выполнить команду:
* Windows:\
`echo $null >> scheduler.db``docker run -p 7540:7540 -v ${PWD}/scheduler.db:/app/scheduler.db finalproject:latest`