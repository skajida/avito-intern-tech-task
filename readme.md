# service-segs

[Условие](docs/task.md) поставленной задачи.

## api

openapi [файлик](api/openapi.yaml) лежит в соседней директории.

## build

Сборку повезло совершать одной командой:

```bash
docker compose -f deployments/docker-compose.yaml up
```

## connect

~~Подключиться можно курлом, но еще~~ есть вариант кидать запросы [постманом](https://www.postman.com/).

```bash
curl -X POST localhost:8080/segs -H 'Content-Type: application/json' -d '{"seg_id": "AVITO_TRAINEE_EXAMPLE_SEGMENT"}'
```

## test

Юнит тестирование, к превеликому сожалению, я сделать не успел. Но мне же не запрещено дописать их в другой ветке?

Немного тестов лежат в этой [папке](test/postman-collections/). Их правда немного, но они минимально показывают работоспособность ручек. Если их можно импортить и разглядывать.

## qa

Действительно, возник один вопрос по ТЗ, конкретно по поводу [третьего дополнительного задания](docs/task.md#доп-задание-3). С минимальным заданием и двумя другими допами должно хватать двух таблиц – таблицы сегментов и таблицы, каким-нибудь способом отражающей взаимосвязь пользователей и сегментов. Вопрос, который возник почти сразу: пользователей чего? Ответ напрашивался: внешнего сервиса.
Решил я эту проблему заглушкой (кажется, это называется мок?) этого сервиса – мне нужно было знать количество пользователей, ну и для приличия сделал метод, который проверяет наличие пользователя в этой foreign базе. Вернее который проверяет, что я не пытаюсь связать свой `user_id` с превышающим количество пользователей айдишником.

## offtop

Я слегка не вписался в дедлайн, хотелось и конфиги покрутить, а не хардкодить друг к другу бд и сервис; до юнит тестов, которые, на мой взгляд, должны неплохо показывать навыки кандидатов я вообще не успел добраться и сожалею, что не очень разумно распределил нагрузку на этой неделе.

Написать достойный отчет я тоже не успеваю :)

Спасибо за это задание! Это был мой первый pet-project на Go и мне понравилось.
