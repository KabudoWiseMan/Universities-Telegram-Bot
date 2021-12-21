# Телеграм бот для подбора университетов

Телеграм бот для подбора университетов с автообновлением базы данных путём парсинга веб-страниц.

Доступен по [ссылке](https://t.me/Choose_University_Bot) или QR-коду:
<img src="https://drive.google.com/uc?export=view&id=1pDZzXoanmPguoXg36kPQr5A1qMYJzam4">

Бот позволяет подобрать подходящие вузы по следующим параметрам:
* Баллы ЕГЭ
* Город
* Профиль и специальность
* Цена обучения
* Наличие вступительных испытаний
* Наличие общежития
* Наличие военной кафедры.

Также есть возможность посмотреть, какие российские вузы находятся в [международном рейтинге QS](https://www.topuniversities.com/qs-world-university-rankings)

Примеры запросов:
<img src="https://drive.google.com/uc?export=view&id=1dmgwYEbDBcW_6xHPQSr1jpNUFcLznbP3">

Использованные инструменты:
* Golang
* Telegram Bot API
* PostgreSQL
* Chromedp.

Развёрнут на Heroku.

Схема базы данных:
<img src="https://drive.google.com/uc?export=view&id=1AmxNtMIEcbmvGeFU1eZe4CRsdOft2EQt">
