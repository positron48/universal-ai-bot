# План реализации AI Telegram-бота для изучения английского

## Архитектура и стек технологий

**Обзор архитектуры:** Бот будет состоять из сервера на Go, работающего на VPS, который обрабатывает сообщения из Telegram и взаимодействует с моделью OpenAI. Telegram будет отправлять обновления боту через webhook на выделенный URL (настроенный через Nginx). Бэкэнд на Go будет минималистичным: без дополнительного кэша или очередей – все данные хранятся локально в SQLite базе. Это упрощает развертывание, устраняя необходимость в внешних сервисах (в соответствии с требованием использовать минимальные зависимости). Каждый Telegram-пользователь рассматривается как отдельный пользователь бота, что будет отражено в базе данных (по уникальному `user_id` Telegram).

**Компоненты решения:**

- **Telegram Bot API (вебхук):** Для приёма сообщений от пользователей. Бот использует Telegram Bot API через вебхук, что обеспечивает мгновенную доставку обновлений без постоянного опроса. Вебхук настраивается на HTTPS URL вашего сервера; Telegram при каждом сообщении пользователя делает POST-запрос на этот адрес с данными обновления. Nginx будет выступать как прокси, принимая HTTPS-запросы от Telegram и перенаправляя их на локальный порт, где слушает Go-приложение.
- **Go (Go 1.x):** Язык бэкенда, обеспечивающий высокую производительность и простую реализацию конкурентной обработки сообщений. Мы выбрали Go, так как он позволяет собрать статически связанный бинарник, удобный для деплоя на VPS, и имеет хорошие библиотеки для Telegram Bot API.
- **Библиотека для Telegram API:** Будем использовать официальную или популярную библиотеку для Telegram ботов на Go. Например, библиотека `go-telegram/bot` предоставляет обёртку без внешних зависимостей (zero-dependencies), что соответствует нашим требованиям[github.com](https://github.com/go-telegram/bot#:~:text=,2 from August 15%2C 2025). Она поддерживает актуальную версию Bot API и умеет работать как через long polling, так и через вебхуки. С помощью этой библиотеки мы легко отправляем и получаем сообщения, обрабатываем команды и колбэки.
- **SQLite база данных:** Для хранения данных используем файл базы SQLite. SQLite не требует развертывания отдельного сервера БД и отлично подходит для небольших проектов на одном сервере. В базе будем хранить информацию о пользователях, словарные карточки и прогресс повторения слов. **Важно:** SQLite позволяет конкурентное чтение, но не поддерживает одновременную запись из нескольких потоков[detunized.net](https://detunized.net/posts/2019-04-01-telegram-bot-in-go-concurrent-sqlite/#:~:text=From the FAQ%3A), поэтому следует либо выполнять операции записи последовательно, либо использовать пул соединений. В рамках нашего бота нагрузка невелика, поэтому можно обойтись одним подключением/потоком записи или небольшим пулом.
- **OpenAI API (через OpenRouter):** Для реализации интеллектуальных функций (проверка грамматики, перевод) подключим API OpenAI. Планируется использовать OpenRouter – это прослойка, совместимая с OpenAI API, дающая единый интерфейс к разным моделям. Мы сможем вызывать, к примеру, модель GPT-3.5 или GPT-4 через OpenRouter, используя OpenAI-протокол. Это потребует указать базовый URL `https://openrouter.ai/api/v1` и передать API-ключ OpenRouter как Bearer-токен в заголовке запросов[docs.langwatch.ai](https://docs.langwatch.ai/integration/go/integrations/openrouter#:~:text=oaioption.WithBaseURL(). Ключ OpenAI/OpenRouter будет храниться в файле окружения `.env` и загружаться при старте приложения, чтобы не хардкодить секреты в коде[go-telegram-bot-api.dev](https://go-telegram-bot-api.dev/#:~:text=Instead of typing the API,receive messages from your Bot). OpenRouter позволяет работать через стандартные SDK OpenAI, задав переменную окружения и базовый URL, поскольку схема запросов почти идентична OpenAI Chat API[openrouter.ai](https://openrouter.ai/docs/api-reference/authentication#:~:text=a Bearer token with your,API key).
- **ENV-конфигурация:** Все чувствительные данные (токен Telegram Bot, API-ключ OpenRouter, и пр.) будут храниться в переменных окружения (например, в `.env` файле) и загружаться при запуске. Такой подход предотвратит утечку секретов в репозиторий и упростит смену ключей[go-telegram-bot-api.dev](https://go-telegram-bot-api.dev/#:~:text=Instead of typing the API,receive messages from your Bot). В `.env` файле ожидаются переменные, например: `TELEGRAM_TOKEN` (токен бота), `OPENROUTER_API_KEY` (ключ OpenRouter), и др. Пакет `os` из стандартной библиотеки или небольшая утилита вроде `godotenv` помогут загрузить эти переменные.

## Пошаговый план реализации

1. **Инициализация проекта и настройка окружения:** Создайте новый модуль Go для бота (командой `go mod init`). Установите необходимые библиотеки: например, пакет для Telegram Bot API (как упомянуто, `go-telegram/bot` с нулевыми зависимостями) и драйвер SQLite (например, `github.com/mattn/go-sqlite3`). В репозитории добавьте файл конфигурации `.env.example` с перечислением требуемых переменных (Telegram token, OpenRouter API key), чтобы было ясно, какие настройки нужны. Реализуйте чтение конфигурации: в функции `main()` загрузите переменные окружения и проверьте, что все необходимые значения установлены. **Важно:** не храните токены в коде – используйте `os.Getenv("TELEGRAM_TOKEN")` и аналогично для OpenAI API ключа[go-telegram-bot-api.dev](https://go-telegram-bot-api.dev/#:~:text=Instead of typing the API,receive messages from your Bot). На этом этапе можно инициализировать логгер для отладки.
2. **Настройка Telegram-бота и вебхука:** Зарегистрируйте бота у @BotFather и получите токен (если ещё не сделано). В коде создайте экземпляр бота с помощью библиотеки, передав токен. Например, для `go-telegram/bot`: `bot.New(token, options...)`. Убедитесь, что бот успешно подключился (метод `GetMe` обычно вызывается автоматически библиотекой для проверки токена). Далее, настройте webhook: необходимо сообщить Telegram, по какому URL отправлять обновления. Это делается либо вызовом метода Bot API `setWebhook` вручную, либо через метод библиотеки. Например, можно вызвать `bot.Methods.SetWebhook(ctx, params)` или использовать готовый метод, если предлагает библиотека[github.com](https://github.com/go-telegram/bot#:~:text=opts %3A%3D []bot.Option). Параметры включают URL вашего сервера (например, `https://yourdomain/bot<token>` или другой путь) и опционально секретный токен. **Настройка Nginx:** На VPS настройте Nginx как обратный прокси: он должен принимать HTTPS запросы к указанному в вебхуке пути и проксировать их на локально запущенный Go-сервис. Например, location `/bot` на Nginx может перенаправлять на `http://127.0.0.1:3000/` (порт, где слушает приложение). Удостоверьтесь, что сертификат SSL настроен (например, с Let’s Encrypt), так как Telegram требует HTTPS для вебхуков. В коде бота запустите HTTP-сервер, который слушает указанный путь. С библиотекой это просто: достаточно вызвать `bot.StartWebhook(...)` и передать `bot.WebhookHandler()` в `http.ListenAndServe`[github.com](https://github.com/go-telegram/bot#:~:text=If you want to use,HTTP handler for your server)[github.com](https://github.com/go-telegram/bot#:~:text=go b) – библиотека сама обработает входящие обновления. Если вебхук секретный токен установлен, убедитесь, что библиотека/код проверяет заголовок `X-Telegram-Bot-Api-Secret-Token`. После запуска сервера Telegram начнёт присылать сообщения бота в ваш обработчик. Для отладки и тестов на этапе разработки, когда HTTPS может быть недоступен, можно временно использовать long polling (метод `GetUpdates`) или инструмент типа `ngrok` для туннелирования webhook’ов.
3. **Проектирование и создание базы данных (SQLite):** Спланируйте схему БД, отражающую основные сущности: **Users**, **Words**, **UserWords (progress)**. Создайте таблицу пользователей, чтобы регистрировать каждого нового пользователя Telegram (поля: `telegram_id` – PRIMARY KEY, имя/username, дата первого использования и т.п.). Каждому аккаунту Telegram соответствует отдельная запись[habr.com](https://habr.com/ru/articles/907716/#:~:text=Первой сложностью стало предложение ChatGPT,фичу отслеживания прогресса изучения слов). Таблица слов будет содержать английские слова или фразы для изучения (поля: `word_id`, `english_word`, `translation_ru` и возможно транскрипция, пример предложения и др. колонки для расширения контента). Можно заранее наполнить эту таблицу набором слов (например, наиболее частотные слова) или оставить пустой, позволяя пользователю пополнять свой словарь. Для хранения прогресса изучения нужна связующая таблица, например **user_word_progress**, где каждой паре (user_id, word_id) ставится статус: выучено или нет, счетчики правильных ответов, дата следующего повторения и т.д. Поля `last_reviewed` и `next_review` помогут реализовать интервальные повторения. **Initial migration:** Реализуйте инициализацию БД при старте приложения: подключитесь к файлу SQLite (например, `database.sqlite` в рабочей директории) с помощью пакета `database/sql`. Выполните SQL DDL для создания таблиц, если они ещё не существуют (команды `CREATE TABLE IF NOT EXISTS ...`). Это можно сделать вручную через `Exec` на старте. Хранение прогресса в БД критично для отслеживания выученных слов[habr.com](https://habr.com/ru/articles/907716/#:~:text=Первой сложностью стало предложение ChatGPT,фичу отслеживания прогресса изучения слов) – без базы данных бот не сможет помнить, какие слова уже показывались и как пользователь справился. Учтите, что SQLite не предназначен для больших нагрузок, но для небольшого бота его производительности хватит. В будущем, если пользователей станет много, можно будет перейти на PostgreSQL или другую СУБД (см. раздел улучшений)[habr.com](https://habr.com/ru/articles/907716/#:~:text=Выбрал неоптимальную базу данных,ИИ конвертировать SQLite в PostgreSQL).
4. **Реализация функционала карточек слов (flashcards):** Добавьте в бот команду или режим для изучения слов, например `/learn` или `/cards`, которая запускает сессию флеш-карт. Алгоритм работы в упрощённом виде: бот выбирает слово из словаря, которое пользователь ещё не выучил (или которое пора повторить согласно расписанию), и отправляет его пользователю. Вместе со словом можно отправить подсказку или пример использования (если есть в базе). Перевод сразу можно **не** показывать, чтобы пользователь попытался вспомнить значение. Пользователь отвечает либо переводом/значением, либо может попросить подсказку/перевод (для этого реализуйте команду `/translate` или кнопку «Показать перевод»). После ответа пользователь отмечает, знал он слово или нет. Это можно реализовать с помощью кнопок: например, прикрепить к сообщению две inline-кнопки – **«Следующее слово»** и **«Отметить как выученное»**, как это делают аналогичные боты[habr.com](https://habr.com/ru/articles/907716/#:~:text=,п). При нажатии «выучено» бот помечает в базе, что слово изучено (например, ставит `learned=true` для записи в таблице прогресса) и больше не будет предлагать его в обычных повторах (либо будет предлагать реже, в зависимости от алгоритма повторения). Кнопка «следующее слово» просто вытаскивает следующий элемент. Если слово помечено выученным, можно сразу показать перевод и пример, чтобы закрепить знание. Реализуйте подсчет прогресса: например, команда `/progress` может показывать, сколько слов выучено из общего списка, или процент прогресса[habr.com](https://habr.com/ru/articles/907716/#:~:text=,«Отметить как выученное»). Для логики повторения используйте принцип ** spaced repetition ** (интервальные повторения). Более сложные алгоритмы, такие как Leitner или SM2, можно добавить позже, но изначально можно задать фиксированные интервалы повторения (например, повторить новое слово через день, через 3 дня, через неделю и т.д.). С помощью полей `next_review` бот может выбирать для каждой сессии слова, у которых наступила дата повторения. Таким образом, обеспечивается *“запоминание/повторение слов”* – слова не теряются, а периодически повторяются, причём более трудные или недавно выученные – чаще, а старые и лёгкие – реже, что соответствует эффективной технике интервальных повторений[en.wikipedia.org](https://en.wikipedia.org/wiki/Spaced_repetition#:~:text=Spaced repetition is an evidence,1). При каждой повторной карточке бот может спрашивать перевод или использование слова, и опять же фиксировать, вспомнил пользователь или нет, обновляя `next_review` (если ошибка – показывать слово скорее, если успех – увеличивать интервал до следующего раза).
5. **Реализация функции перевода (ENG↔RU):** Предусмотрите возможность перевода слов и фраз как с английского на русский, так и с русского на английский. Проще всего сделать единую команду `/translate`, за которой указывается текст. Бот должен определить язык входящего текста (либо указать два разных команды, например `/en2ru` и `/ru2en`). Можно автоматизировать: если текст содержит кириллицу, считать его русским и переводить на английский, иначе – переводить на русский. Для перевода используйте вызов OpenAI API через OpenRouter. Например, при запросе перевода английского текста на русский бот формирует prompt для модели GPT-3.5: *«Translate this text to Russian: '<USER_TEXT>'»*. Модель вернёт перевод, бот отправляет его пользователю. В обратном направлении аналогично (prompt: *«Переведи на английский: '<USER_TEXT>'»*). Такой подход использует мощность AI для качественного перевода, превосходящего по качеству простые словарные подстановки. В коде целесообразно инкапсулировать эту логику в отдельной функции `translate(text, targetLang)`. При её вызове формируется HTTP-запрос к эндпоинту OpenRouter (например, `POST /api/v1/chat/completions`), с авторизационным заголовком `Authorization: Bearer <OPENROUTER_API_KEY>` и JSON-телом запроса (модель, сообщения и т.д.). **Пример:** модель `openai/gpt-3.5-turbo` прекрасно справляется с переводами. API-ключ OpenRouter берётся из переменной окружения, как настроено ранее[docs.langwatch.ai](https://docs.langwatch.ai/integration/go/integrations/openrouter#:~:text=oaioption.WithBaseURL(). Библиотека OpenAI для Go (если используется) настраивается на базовый URL OpenRouter и будет использовать наш ключ как OpenAI-ключ[docs.langwatch.ai](https://docs.langwatch.ai/integration/go/integrations/openrouter#:~:text=oaioption.WithBaseURL(). После получения ответа от модели, бот отправляет переведенный текст в чат. Учтите ограничения: большие тексты переводить дольше и дороже, поэтому можно ограничить длину запроса или уведомлять пользователя, если текст слишком большой. Также, оберните вызов API в обработку ошибок – если сеть или API недоступны, бот должен ответить пользователю чем-то вроде “⚠️ Извините, не удалось перевести, попробуйте позже”.
6. **Реализация функции проверки английского (грамматика и ошибки):** Добавьте команду, например `/check` или `/grammar`, которая будет принимать от пользователя английское предложение или абзац для проверки. Бот, получив такой запрос, с помощью AI-подсистемы проанализирует текст и найдёт ошибки. Реализация: сформируйте запрос к модели GPT-3.5/GPT-4, где в prompt описано: *«Find grammar and spelling mistakes in the following English text and suggest corrections. Text: '<TEXT>'.»*. Можно попросить модель вернуть отредактированное предложение или указать ошибки. На начальном этапе достаточно вернуть исправленное предложение/текст. Например, пользователь отправил: "*I has a apple.*" – бот ответит: "*I **have** an apple.*". Чтобы сделать сервис более обучающим, можно немного усложнить: пусть модель не только исправит, но и **пояснит** ошибки. Тогда ответ модели может содержать объяснение (например: "*Ошибка: неправильная форма глагола 'have'. Правильно: 'I have an apple'.*"). Бот отправит пользователю оба варианта: и исправленный текст, и комментарий. (При оформлении ответа можно использовать Markdown-разметку, чтобы выделять исправления жирным или зачёркнутым текстом для наглядности). Внедрение этой функции с помощью OpenAI очень похоже на перевод: меняется только формулировка задачи для модели. Также можно задать модели роль (system message), например: *“You are an English teacher bot...”*, чтобы она отвечала более педагогично. Обработайте ошибки API: если модель не смогла проанализировать (что маловероятно), верните пользователю уведомление об ошибке. Данная функция реализует *“проверку английского”*: пользователь получает мгновенную обратную связь на свой текст. Это особенно полезно для самостоятельного письма на английском.
7. **Тестирование и деплой на сервер:** После реализации основных функций, тщательно протестируйте бота локально. Напишите простые юнит-тесты для отдельных функций (например, парсер команд, корректность формирования запросов к API, функции работы с БД). Затем проведите интеграционное тестирование: запустите бота и отправьте ему сообщения в Telegram, проверяя все команды: добавление пользователя через `/start`, цикл карточек, команды перевода и проверки грамматики. Убедитесь, что бот корректно обновляет базу (помечает слова как выученные, планирует повторения). **Логирование:** добавьте логирование важных событий (запуск, получение сообщения, вызов API, ошибки БД) – это поможет при отладке на сервере. Настройте запуск на VPS: скомпилируйте Go-приложение (`CGO_ENABLED=1` может понадобиться для SQLite драйвера) и скопируйте бинарник и файл базы (если нужен) на сервер. Настройте systemd-сервис или Docker-контейнер для постоянной работы бота. Убедитесь, что в окружении на сервере заданы необходимые переменные (можно использовать файл `.env` и загрузить через systemd unit). **Настройка вебхука в продакшене:** ещё раз выполните `setWebhook` с вашим продоменом, если до этого тестировали локально иначе. Проверьте через Bot API `getWebhookInfo`, что вебхук установлен и нет ошибок. На Nginx убедитесь, что прокси прокидывает нужные заголовки (например, `Host` и т.д.) и что секретный токен (если задан) совпадает. **Безопасность и лучшие практики:** регулярно обновляйте токены/ключи, не публикуйте их. Следите за обновлениями API Telegram и OpenAI. По возможности, ограничьте доступ бота к интернету кроме нужных API, и настройте брандмауэр на VPS. Бот не использует внешние БД или кеши, что упрощает поддержку: резервное копирование SQLite-файла будет достаточным для сохранности данных. При росте нагрузки обратите внимание на потенциальные узкие места: последовательная обработка сообщений (Go легко позволяет обрабатывать параллельно, но тогда учтите блокировки при доступе к SQLite) и задержки от OpenAI API (можно внедрить асинхронную обработку или очередь, если ответы AI станут медленными, но это уже усложнение сверх минимального решения).

## Возможные доработки и улучшения

- **Диалоговый режим с AI-наставником:** Помимо командного режима, можно реализовать свободное общение. В таком режиме бот выступает в роли собеседника на английском: пользователь переписывается на языке, а бот поддерживает разговор, одновременно корректируя ошибки. При этом бот продолжает диалог, но если пользователь допускает ошибку, бот тактично исправляет и объясняет ее перед тем, как ответить по теме разговора. Это требует тонкой настройки промптов для модели (роль учителя), чтобы она умела одновременно быть собеседником и корректором. Такой интерактивный режим повысит вовлечённость и поможет учиться на собственных ошибках, не прерывая беседу.
- **Улучшенный UX с кнопками и меню:** Чтобы упростить взаимодействие, можно использовать возможности Telegram UI. Например, реализовать кастомную клавиатуру для основных функций (кнопки: “Новая карточка”, “Мой прогресс”, “Перевести”, “Проверить текст”). В сессиях изучения слов стоит использовать inline-кнопки для ответа вместо ввода команд – как упоминалось, “Следующее слово”, “Выучено” и “Показать перевод”. Также можно добавить кнопки навигации: «Назад» для возврата в главное меню и переключатели режима. Эти улучшения были предложены самим автором бота при его разработке[habr.com](https://habr.com/ru/articles/907716/#:~:text=Примеры запросов%2C которые я использовал%3A) и значительно улучшают удобство использования.
- **Мультиязычный интерфейс:** Добавьте поддержку русского и английского интерфейса. Пользователь сможет выбирать язык общения с ботом (например, через команду `/language` или кнопку). В зависимости от выбора, бот будет присылать инструкции, кнопки и сообщения либо на русском, либо на английском. Это полезно, чтобы начинающим было проще (русский интерфейс), а продвинутым – полное погружение в английский. Переключение языков можно хранить в профиле пользователя (в БД) и применять ко всем исходящим сообщениям бота[habr.com](https://habr.com/ru/articles/907716/#:~:text=Примеры запросов%2C которые я использовал%3A).
- **Расширение контента карточек:** В текущей версии карточки могут содержать только слово и перевод. Можно расширить их информацией: добавьте транскрипцию (произношение), один-два примера использования слова в предложении и синонимы/антонимы. Эти данные можно получить автоматизированно. Например, используя сторонний словарный API либо сгенерировать через всё тот же OpenAI. Автор похожего проекта упоминает, что сгенерировал базу из 3000 слов с транскрипциями, переводами и примерами при помощи AI-сервиса (DeepSeek) вместо ручного заполнения[habr.com](https://habr.com/ru/articles/907716/#:~:text=Писать вручную 3000 слов с,слова из PDF в Excel). Вы тоже можете реализовать команду для обогащения слов: при добавлении нового слова бот сам обращается к AI, чтобы получить перевод и пример предложения, и сохраняет их. Это минимизирует ручной труд по наполнению словаря.
- **Продвинутая система повторений:** Внедрение полноценного алгоритма *spaced repetition* (интервальных повторений) для планирования карточек. Например, реализовать алгоритм SuperMemo 2 (SM-2), который используется в Anki. Этот алгоритм рассчитывает интервалы на основе оценки ответа пользователя (как хорошо он помнит слово). Правильно отвеченные слова будут появляться через увеличивающиеся промежутки (2 дня, 6 дней, 14 дней, ...), а ошибки сбрасывают интервал. Такое расписание оптимизирует запоминание и доказано повышает эффективность обучения[en.wikipedia.org](https://en.wikipedia.org/wiki/Spaced_repetition#:~:text=Spaced repetition is an evidence,1). Потребуется хранить в `user_word_progress` коэффициенты сложности карточки и текущий интервал. Реализация подобного алгоритма усложнит логику, но значительно повысит образовательную ценность бота.
- **Масштабирование и производительность:** Если число пользователей и объем данных вырастет, рассмотрите переход от SQLite к серверной СУБД (PostgreSQL, MySQL). В случае деплоя на нескольких нодах или в облаке файл SQLite неудобен, что отмечалось в опыте создания аналогичного бота[habr.com](https://habr.com/ru/articles/907716/#:~:text=Выбрал неоптимальную базу данных,ИИ конвертировать SQLite в PostgreSQL). PostgreSQL более надёжен для конкурентных записей и больших данных. Кроме того, можно внедрить кэширование результатов частых запросов (например, кеш переводов для популярных фраз) в память или с помощью Redis – однако помните, что это увеличит сложность развертывания.
- **Дополнительные функции:** В будущем можно добавить проверки произношения (при помощи распознавания речи: пользователь отправляет голосовое сообщение, бот конвертирует его в текст и проверяет, либо синтезирует эталонное произношение). Также полезной может быть регулярная рассылка “слова дня” или подборки новых слов по расписанию – для этого бот может иметь шедулер (cron-задачу) или запускаться периодически через планировщик, отправляя контент пользователям, которые подписались на рассылку[vc.ru](https://vc.ru/education/2172922-telegram-bot-pocket-dictionary-umnaya-rassylka#:~:text=Telegram,частью речи%3B Живой пример). Такие проактивные рассылки удерживают интерес пользователей.

Проектируя и улучшая бота, старайтесь следовать **best practices** разработки: держите код аккуратным и раздельным (логика бота, работа с БД, интеграция с AI – в отдельных модулях), проверяйте ошибки на каждом шаге, не допускайте блокировки потока обработки (вебхук-хендлер должен быстро отвечать Telegram, а тяжёлые задачи можно выполнять асинхронно). Документируйте код и обновляйте документацию по мере добавления новых фич. Таким образом, у вас получится поддерживаемый AI-бот, который поможет пользователям эффективно изучать английский язык в удобном формате чат-мессенджера.



## Технический пошаговый план (от Git до продакшна)

### Шаг 0. Инициализация репозитория и проекта

1. **Создать репозиторий и базовые файлы**
2. **.gitignore и пример окружения**
3. **Структура каталогов (best practices, без перегруза)**

### Шаг 1. Зависимости

Минимальный набор:

```
go get github.com/go-telegram/bot@latest
go get github.com/mattn/go-sqlite3@latest     # CGO, стабильный
# или: go get modernc.org/sqlite@latest       # без CGO, тяжелее по модулям
# опционально:
go get github.com/joho/godotenv@latest
```

### Шаг 2. Конфигурация и запуск

`internal/config/config.go` — загрузка env, валидация, дефолты. `.env` не коммитим, держим `.env.example`.

### Шаг 3. Схема БД и миграции (SQLite)

В MVP проще хранить **карточки на пользователя** в одной таблице (не общий словарь), чтобы AI‑генерация шла напрямую per‑user.

`migrations/001_init.sql`:

```
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA busy_timeout = 5000;

/* Пользователи */
CREATE TABLE IF NOT EXISTS users (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  telegram_id  INTEGER NOT NULL UNIQUE,
  username     TEXT,
  lang         TEXT NOT NULL DEFAULT 'ru',
  created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

/* Глобальные карточки */
CREATE TABLE IF NOT EXISTS cards (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  en_text      TEXT NOT NULL,              -- слово/фраза
  ru_text      TEXT,                       -- базовый перевод
  ipa          TEXT,                       -- транскрипция
  pos          TEXT,                       -- часть речи
  example_en   TEXT,
  example_ru   TEXT,
  topic        TEXT,                       -- свободная метка темы
  level        TEXT,                       -- A1..C2 (опц.)
  source       TEXT NOT NULL DEFAULT 'ai', -- 'ai'|'manual'
  created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

/* Нормализация для антидубликатов: LOWER/TRIM без лишних пробелов */
CREATE TABLE IF NOT EXISTS cards_norm (
  card_id      INTEGER PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
  norm_en      TEXT NOT NULL,              -- en_text нормализованный
  norm_pos     TEXT,                       -- pos нормализованный
  UNIQUE (norm_en, norm_pos)
);

/* Персональный прогресс пользователя по карточке */
CREATE TABLE IF NOT EXISTS user_cards (
  user_id       INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  card_id       INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
  status        TEXT NOT NULL DEFAULT 'active',   -- active|suspended|leeched
  ef            REAL NOT NULL DEFAULT 2.5,
  repetitions   INTEGER NOT NULL DEFAULT 0,
  interval_days INTEGER NOT NULL DEFAULT 0,
  last_review   DATETIME,
  next_review   DATETIME,
  difficulty    INTEGER NOT NULL DEFAULT 0,       -- 0..5
  added_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, card_id)
);

CREATE TABLE IF NOT EXISTS sessions (
  user_id     INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  state       TEXT NOT NULL,
  payload     TEXT,
  updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

/* Индексы под выборки */
CREATE INDEX IF NOT EXISTS idx_user_cards_due
  ON user_cards(user_id, next_review);
CREATE INDEX IF NOT EXISTS idx_cards_topic
  ON cards(topic);
CREATE INDEX IF NOT EXISTS idx_cards_level
  ON cards(level);
```

Инициализация БД: при старте включить WAL и таймаут:

```
PRAGMA journal_mode = WAL;
PRAGMA busy_timeout = 5000;
```

### Шаг 4. I18N слой (RU/EN) + переключение языка

Минимальный словарь‑ключи в коде:

```
// internal/i18n/messages.go
package i18n

type Lang string
const (
    RU Lang = "ru"
    EN Lang = "en"
)

var dict = map[Lang]map[string]string{
    RU: {
        "start": "Привет! Я помогу учить английский. Выберите режим:",
        "menu.learn": "🃏 Карточки",
        "menu.translate": "🔁 Перевод",
        "menu.check": "🧪 Проверка текста",
        "menu.lang": "🌐 Язык интерфейса",
        "lang.choose": "Выберите язык интерфейса",
        "lang.set.ru": "Интерфейс переключен на русский.",
        "lang.set.en": "Interface switched to English.",
        "cards.none_due": "На сегодня карточек к повторению нет. Сгенерировать новые?",
        "cards.next": "Следующая карточка:",
        "cards.mark.known": "✅ Знал(а)",
        "cards.mark.unknown": "❌ Не знал(а)",
        "cards.show.answer": "Показать ответ",
        "ai.generating": "Генерирую новые карточки…",
        // ...
    },
    EN: {
        "start": "Hi! I’ll help you learn English. Choose a mode:",
        "menu.learn": "🃏 Flashcards",
        "menu.translate": "🔁 Translate",
        "menu.check": "🧪 Check text",
        "menu.lang": "🌐 Interface language",
        "lang.choose": "Choose interface language",
        "lang.set.ru": "Интерфейс переключен на русский.",
        "lang.set.en": "Interface switched to English.",
        "cards.none_due": "No cards due today. Generate new ones?",
        "cards.next": "Next card:",
        "cards.mark.known": "✅ I knew it",
        "cards.mark.unknown": "❌ I didn’t know",
        "cards.show.answer": "Show answer",
        "ai.generating": "Generating new cards…",
        // ...
    },
}

func T(lang Lang, key string) string {
    if m, ok := dict[lang]; ok {
        if s, ok := m[key]; ok { return s }
    }
    return key
}
```

Команда `/lang` показывает inline‑клавиатуру RU/EN, по нажатию мы сохраняем язык в `users.lang`.

### Шаг 5. Интервальные повторения (SM‑2) — MVP

Упрощённая SM‑2 (качество `q ∈ [0..5]`; для “знал” ставим 4–5, “не знал” — 2–3):

```
// internal/srs/sm2.go
package srs

import "math"

type State struct {
    EF          float64 // easiness factor
    Repetitions int
    Interval    int // days
}

func Review(s State, quality int) State {
    if quality < 0 { quality = 0 }
    if quality > 5 { quality = 5 }

    if quality < 3 {
        s.Repetitions = 0
        s.Interval = 1
    } else {
        s.Repetitions++
        switch s.Repetitions {
        case 1:
            s.Interval = 1
        case 2:
            s.Interval = 6
        default:
            s.Interval = int(math.Round(float64(s.Interval) * s.EF))
        }
    }

    s.EF = s.EF + (0.1 - float64(5-quality)*(0.08 + float64(5-quality)*0.02))
    if s.EF < 1.3 {
        s.EF = 1.3
    }
    return s
}
```

В БД обновляем `ef, repetitions, interval_days, last_review, next_review`.

### Шаг 6. Интеграция с OpenRouter (OpenAI‑совместимо) — без SDK

Минимальная зависимость: чистый `net/http`.

```
// internal/ai/client.go
package ai

import (
  "bytes"
  "context"
  "encoding/json"
  "net/http"
  "time"
)

type Client struct {
    baseURL string
    apiKey  string
    http    *http.Client
    model   string // напр. "openai/gpt-4o-mini" или "openai/gpt-3.5-turbo"
}

func New(baseURL, apiKey, model string) *Client {
    return &Client{
        baseURL: baseURL,
        apiKey:  apiKey,
        model:   model,
        http: &http.Client{ Timeout: 25 * time.Second },
    }
}

type ChatReq struct {
    Model    string      `json:"model"`
    Messages []ChatMsg   `json:"messages"`
    Temperature float64  `json:"temperature,omitempty"`
}

type ChatMsg struct {
    Role    string `json:"role"`   // "system" | "user" | "assistant"
    Content string `json:"content"`
}

type Choice struct {
    Message ChatMsg `json:"message"`
}
type ChatResp struct {
    Choices []Choice `json:"choices"`
}

func (c *Client) Chat(ctx context.Context, msgs []ChatMsg) (string, error) {
    reqBody := ChatReq{ Model: c.model, Messages: msgs }
    b, _ := json.Marshal(reqBody)

    httpReq, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewReader(b))
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
    // (опционально для OpenRouter) httpReq.Header.Set("HTTP-Referer", "https://bot.example.com")
    // httpReq.Header.Set("X-Title", "Telegram English Bot")

    resp, err := c.http.Do(httpReq)
    if err != nil { return "", err }
    defer resp.Body.Close()

    var cr ChatResp
    if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil { return "", err }
    if len(cr.Choices) == 0 { return "", nil }
    return cr.Choices[0].Message.Content, nil
}
```

#### AI‑генерация карточек (JSON‑схема на выходе)

Промпт к модели (RU или EN — неважно, лучше EN для стабильности):

```
const CardGenSystem = `You are a vocabulary builder for Russian learners of English.
Return STRICT JSON only.`
func CardGenUserPrompt(topic string, count int, level string) string {
  if topic == "" { topic = "general everyday English" }
  if level == "" { level = "A2" }
  return fmt.Sprintf(`Generate %d English vocabulary items for topic "%s" at level %s.
Return ONLY a JSON array of objects with keys:
"en", "ru", "ipa", "pos", "example_en", "example_ru".
`, count, topic, level)
}
```

Парсим JSON, валидируем поля и вставляем в `cards`.

#### Перевод и проверка грамматики

- **Перевод**: короткий системный промпт “You translate between English and Russian.”, детект языка по наличию кириллицы.
- **Проверка** (`/check`): системный промпт “You are an English teacher… Explain errors briefly.”; просим вернуть: исправленный текст + короткие пояснения.

### Шаг 7. Telegram‑бот, вебхук и маршрутизация

Используем `github.com/go-telegram/bot` (нулевые внешние зависимости).

```
// cmd/bot/main.go
func main() {
  cfg := config.Load()                 // env
  db := repo.MustOpen(cfg.DBPath)      // инициализация + PRAGMA + миграции
  aiClient := ai.New(cfg.OpenRouterBaseURL, cfg.OpenRouterAPIKey, cfg.Model)
  tgb := telegram.MustNew(cfg.TelegramToken, cfg.WebhookSecret, cfg.BindAddr, cfg.ExternalURL)

  // регистрация хендлеров:
  tgb.HandleCommand("start", handlers.Start(db))
  tgb.HandleCommand("lang",  handlers.Lang(db))
  tgb.HandleCommand("cards", handlers.Cards(db, aiClient))
  tgb.HandleCommand("review",handlers.Review(db))       // вызвать due‑карточки
  tgb.HandleCommand("translate", handlers.Translate(aiClient, db))
  tgb.HandleCommand("check", handlers.Check(aiClient, db))
  tgb.HandleDefault(handlers.Fallback(db))              // свободный текст по контексту

  tgb.RunWebhook() // http.ListenAndServe + setWebhook + обработка
}
```

Inline‑кнопки: **Знал / Не знал / Показать ответ / Следующая**. Состояние сессии — в `sessions`.

Потоки: обработчики должны быстро отвечать Telegram (200 ОК). Вызовы AI делаем синхронно, но с таймаутом и “печатает…” (sendChatAction), чтобы UX был живым.

### Шаг 8. Поведение MVP “Карточки”

1. `/cards`:
   - Выбрать **due** карточку: `WHERE user_id=? AND (next_review IS NULL OR next_review<=NOW()) ORDER BY next_review NULLS FIRST LIMIT 1`.
   - Если нет — предложить сгенерировать **N новых** (напр. 10 штук). При согласии: показать `i18n.ai.generating`, дернуть AI, распарсить JSON, вставить в БД, показать первую карточку.
2. Сообщение карточки (скрываем перевод): текст + кнопки.
3. По нажатию:
   - **Показать ответ** → редактируем сообщение, добавляя перевод/пример.
   - **Знал/Не знал** → применяем SM‑2, обновляем поля и `next_review = now + interval_days`.
   - **Следующая** → выбираем следующую due.

### Шаг 9. Тесты (ключевые)

- `srs` — покрыть `Review` для крайних значений `q`.
- `ai` — тест парсинга JSON карточек (table‑driven), защита от мусора.
- `repo` — интеграционный тест CRUD на временной SQLite.

### Шаг 10. Сборка, Nginx, TLS, systemd

**Сборка:**

```
make build
# Makefile:
# build:
#   CGO_ENABLED=1 go build -o bin/tg-english-bot ./cmd/bot
```

**Nginx (прокси вебхука)** — `/etc/nginx/sites-available/bot.conf`:

```
server {
  server_name bot.example.com;
  listen 443 ssl http2;
  ssl_certificate     /etc/letsencrypt/live/bot.example.com/fullchain.pem;
  ssl_certificate_key /etc/letsencrypt/live/bot.example.com/privkey.pem;

  location /tg/webhook {
    proxy_pass         http://127.0.0.1:3000;
    proxy_set_header   Host $host;
    proxy_set_header   X-Forwarded-For $remote_addr;
    proxy_set_header   X-Forwarded-Proto https;
    proxy_http_version 1.1;
    client_max_body_size 5m;
  }
}
```

(HTTP → 301 на HTTPS). TLS — через `certbot --nginx`.

**systemd‑юнит** `/etc/systemd/system/tg-english-bot.service`:

```
[Unit]
Description=Telegram English Bot
After=network.target

[Service]
User=bot
WorkingDirectory=/opt/tg-english-bot
EnvironmentFile=/opt/tg-english-bot/.env
ExecStart=/opt/tg-english-bot/bin/tg-english-bot
Restart=on-failure
RestartSec=5
NoNewPrivileges=true
ProtectSystem=full
ProtectHome=true
PrivateTmp=true
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
```

Дальше:

```
sudo systemctl daemon-reload
sudo systemctl enable --now tg-english-bot
```

**Webhook установка** — при старте бот сам вызывает `setWebhook` с `WEBHOOK_EXTERNAL_URL` и `WEBHOOK_SECRET`. Проверить:

- `getWebhookInfo` в логах/методом — должно быть `url=...` и `has_custom_certificate=false`, `pending_update_count=0`.

**Бэкапы**: cron‑копия `database.sqlite` (WAL‑режим → копировать и `-wal` при необходимости), ежедневная ротация логов.