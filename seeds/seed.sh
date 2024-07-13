#!/bin/bash

echo "Ожидание базы данных..."

until nc -z $POSTGRES_HOST $POSTGRES_PORT; do
    sleep 1
done
echo "База данных доступна"

echo "Запуск скриптов сидов"
export PGPASSWORD=$POSTGRES_PASSWORD

# запуск скриптов сидов для юзерков
psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_NAME -f /seeds/01_users_seed.sql
users_seed_status=$?
if [ $users_seed_status -eq 0 ]; then
    echo "Скрипт users_seed.sql успешно выполнен"
else
    echo "Ошибка при выполнении скрипта users_seed.sql"
    exit 1
fi

# Запуск скриптов сидов для тасок
psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_NAME -f /seeds/02_tasks_seed.sql
tasks_seed_status=$?
if [ $tasks_seed_status -eq 0 ]; then
    echo "Скрипт tasks_seed.sql успешно выполнен"
else
    echo "Ошибка при выполнении скрипта tasks_seed.sql"
    exit 1
fi

echo "Скрипты сидов успешно отработали, база данных пополнена"
