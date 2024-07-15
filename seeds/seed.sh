#!/bin/bash

echo "Waiting for Database"

until nc -z $POSTGRES_HOST $POSTGRES_PORT; do
    sleep 1
done
echo "Database available"

echo "Running seed scripts"
export PGPASSWORD=$POSTGRES_PASSWORD

# запуск скриптов сидов для юзерков
psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_NAME -f /seeds/01_users_seed.sql
users_seed_status=$?
if [ $users_seed_status -eq 0 ]; then
    echo "The users_seed.sql script was executed successfully"
else
    echo "Error when executing the users_seed.sql script"
    exit 1
fi

# Запуск скриптов сидов для тасок
psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_NAME -f /seeds/02_tasks_seed.sql
tasks_seed_status=$?
if [ $tasks_seed_status -eq 0 ]; then
    echo "The tasks_seed.sql script was executed successfully"
else
    echo "Error when executing the tasks_seed.sql script"
    exit 1
fi

echo "The seed scripts have been successfully processed, the database has been replenished"
