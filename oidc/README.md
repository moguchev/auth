# OpenID Connect

Туториал:

1. Открываем в браузере http://localhost:8080/
2. Заходим под `admin`
3. Создаем новый Realm `example`
4. Создаем нового клиента `example-client`:
   1. `Client authentication`: `On`
   2. `Authorization`: `On`
   3. `Valid redirect URIs`: `http://localhost:3000/oidc/callback`
   4. `Login theme`: `keycloak`
5. Выставляем переменные окружения:
   1. `export OIDC_CLIENT_SECRET=<secret>` - секрет можно найти `example-client` → `Credentials` → `Client Secret`
   2. `export OIDC_CLIENT_ID=example-client`
6. Создаем пользователей:
   1. `Users` → `Create new user`
   2. Создаем пользователя `Bob`:
      1. `Email verified`: `true`
      2. `Username`: `bob`
      3. `Email`: `bob@gmail.com`
      4. `First name`: `Bob`
      5. `Last name`: `Bob`
      6. Нажимаем `Create`
      7. `Credentials` → устанавливаем пароль (не временный)
7. Запускаем `auth_service`:
```sh
cd auth_service && go run .
```
8. Запускаем `greeting_service`:
```sh
cd greeting_service && go run .
```
9. Переходим на страницу авторизации: http://localhost:3000/login/oidc/authorize
10. Логинимся под учетной записью `Bob`
11. Сохраняем `access_token`
12. Делаем запрос в `greeting_service`:
```sh
curl http://localhost:8081/hello -H "Authorization: Bearer <access_token>"
```
13. Получаем ответ:
```json
{"message":"Hello Bob Bob!"}
```