# Менеджер паролей GophKeeper

GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю надёжно и безопасно хранить логины, пароли, бинарные данные и прочую приватную информацию.

## [Техническое задание](technical_task.md)

## Запуск

### Сертификат

Чтобы сгенерировать самоподписанный сертификат:

```shell
make cert-gen
```

### Сервер

При необходимости можно пересобрать прото-файлы:

```shell
make proto name=user
make proto name=creds
```

Сначала нужно собрать бинарник:
```shell
make build-server
```

Далее необходимо запустить БД:
```shell
docker-compose up --build
```

Потом в другом терминале запустить сам сервер:
```shell
./bin/cli
```

### Клиент

Сначала нужно собрать бинарник:
```shell
make build-cli
```

Далее можно использовать команды CLI:
```shell
./bin/cli register username password
./bin/cli login username password
./bin/cli logout

./bin/cli add-creds https://example.com/ login password
./bin/cli list-creds
./bin/cli get-creds 1234
./bin/cli delete-creds 1234
```
