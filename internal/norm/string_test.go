package norm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	tests := map[string]string{
		"00 Вступление. Почему мы\u00a0решили записать книгу о страсти":     "00 Вступление. Почему мы решили записать книгу о страсти",
		"01 Глава 1. Мы смотрим, но не всегда видим ":                       "01 Глава 1. Мы смотрим, но не всегда видим",
		"01 Глава 1. Страсть-обращаться с\u00a0осторожностью":               "01 Глава 1. Страсть-обращаться с осторожностью",
		"01 Глава 1. Возражение\u00a0—\u00a0 естественная часть общения":    "01 Глава 1. Возражение —  естественная часть общения",
		"01 Введение. Я владела своей тайной, а моя тайна владела мною": "01 Введение. Я владела своей тайной, а моя тайна владела мною",
		"01 Введение. Разговор начистоту - будем предельно откровенны ":     "01 Введение. Разговор начистоту - будем предельно откровенны",
		"01 Дисклеймер": "01 Дисклеймер",
		"01. Правило № 1. Физическая нагрузка стимулирует работу мозга, часть 1": "01. Правило No 1. Физическая нагрузка стимулирует работу мозга, часть 1",
		"02 Не позволяйте своим талантам умереть вместе с вами ":                 "02 Не позволяйте своим талантам умереть вместе с вами",
		"02 Глава 1. Как понять, какой блог вам нужен":                          "02 Глава 1. Как понять, какой блог вам нужен",
		"02 Глава 2. Анализ - сконцентрированы\u00a0ли вы на самом важном":       "02 Глава 2. Анализ - сконцентрированы ли вы на самом важном",
		"02 Глава 2. Сравним возражение с жалобой и претензией ":                 "02 Глава 2. Сравним возражение с жалобой и претензией",
		"02 Глава 2. Эволюция страсти-от\u00a0страдания к\u00a0любви":            "02 Глава 2. Эволюция страсти-от страдания к любви",
		"02. 1 Делайте больше дел одновременно":                                 "02. 1 Делайте больше дел одновременно",
		"02. Урок №1. Найдите стимул":                                            "02. Урок No1. Найдите стимул",
		"02. Правило № 1. Физическая нагрузка стимулирует работу мозга, часть 2": "02. Правило No 1. Физическая нагрузка стимулирует работу мозга, часть 2",
		"03 Глава 3. Как обнаружить и\u00a0разжечь в\u00a0себе страсть":          "03 Глава 3. Как обнаружить и разжечь в себе страсть",
		"03 Глава 3. Нет, все сразу не\u00a0получится":                           "03 Глава 3. Нет, все сразу не получится",
		"03 Глава 3. Разбитые сердца и\u00a0сломанные ноги":                      "03 Глава 3. Разбитые сердца и сломанные ноги",
		"03 Привычка № 1. Стремиться к пониманию":                                "03 Привычка No 1. Стремиться к пониманию",
		"03 Рабочий мусор": "03 Рабочий мусор",
		"03. 2 Отдавайте энергию и время тем, кому это действительно нужно":                          "03. 2 Отдавайте энергию и время тем, кому это действительно нужно",
		"03. Урок №2. Соберите союзников":                                                              "03. Урок No2. Соберите союзников",
		"03. Правило № 2. Выживание: мозг человека тоже эволюционировал, часть 1":                      "03. Правило No 2. Выживание: мозг человека тоже эволюционировал, часть 1",
		"04 04 Глава 4. Как работать со сложной аудиторией":                                          "04 04 Глава 4. Как работать со сложной аудиторией",
		"04 Глава 2 Создание и обустройство мира":                                                     "04 Глава 2 Создание и обустройство мира",
		"04 Глава 2. Движущийся мир #01":                                                              "04 Глава 2. Движущийся мир #01",
		"04 Глава 4. У\u00a0справедливости вкус\u00a0шоколада":                                         "04 Глава 4. У справедливости вкус шоколада",
		"04 Глава 4. Как наши роли влияют на\u00a0нас":                                                 "04 Глава 4. Как наши роли влияют на нас",
		"04 Глава 4. Когда страсть заводит не\u00a0туда":                                               "04 Глава 4. Когда страсть заводит не туда",
		"04 Глава 4. Высвободите место в\u00a0работе для остальных сторон вашей личности":              "04 Глава 4. Высвободите место в работе для остальных сторон вашей личности",
		"04 Привычка № 2. Генерировать энергию":                                                        "04 Привычка No 2. Генерировать энергию",
		"04. 3 Меняйте жизнь обдуманно. Меняйте жизнь радикально":                                    "04. 3 Меняйте жизнь обдуманно. Меняйте жизнь радикально",
		"04. Миф №1. Иностранный язык можно выучить":                                                   "04. Миф No1. Иностранный язык можно выучить",
		"04. Урок №3. Разберитесь в своих чувствах":                                                    "04. Урок No3. Разберитесь в своих чувствах",
		"04. Глава 3. Подъем на\u00a0лифте настроения. Преимущества":                                   "04. Глава 3. Подъем на лифте настроения. Преимущества",
		"04. Глава 3. Действия. Первые шага для старта":                                               "04. Глава 3. Действия. Первые шага для старта",
		"04. Правило № 2. Выживание: мозг человека тоже эволюционировал, часть 2":                      "04. Правило No 2. Выживание: мозг человека тоже эволюционировал, часть 2",
		"05 День 2. Природный энергетик":                                                              "05 День 2. Природный энергетик",
		"05 Глава 2. Движущийся мир #02":                                                              "05 Глава 2. Движущийся мир #02",
		"05 Глава 5. Мы можем одновременно шагать и жевать жвачку - ничего больше ":                    "05 Глава 5. Мы можем одновременно шагать и жевать жвачку - ничего больше",
		"05 Глава 5. Что мешает быть собой и\u00a0поступать эффективно":                                "05 Глава 5. Что мешает быть собой и поступать эффективно",
		"05 Глава 5. Перед тем как согласиться на\u00a0новую работу, задайте себе три вопроса":         "05 Глава 5. Перед тем как согласиться на новую работу, задайте себе три вопроса",
		"05 Привычка № 3. Обдумывать необходимость":                                                    "05 Привычка No 3. Обдумывать необходимость",
		"05. Миф №2. Выучить язык и начать на нем говорить легче в детстве":                            "05. Миф No2. Выучить язык и начать на нем говорить легче в детстве",
		"05. Урок №4. Долой хлам!":                                                                     "05. Урок No4. Долой хлам!",
		"05. Правило № 3. Мозг каждого человека имеет различную электропроводимость нейронов, часть 1": "05. Правило No 3. Мозг каждого человека имеет различную электропроводимость нейронов, часть 1",
		"05. Принцип №1. Запретный плод сладок":                                                        "05. Принцип No1. Запретный плод сладок",
		"06 День 3. Бедные мои яйца":                                                                  "06 День 3. Бедные мои яйца",
		"06 Глава 2. Движущийся мир #03":                                                              "06 Глава 2. Движущийся мир #03",
		"06 Глава 4 Достойные Вальхаллы - люди-герои":                                                 "06 Глава 4 Достойные Вальхаллы - люди-герои",
		"06 Глава 5. Перестройка":                                                                     "06 Глава 5. Перестройка",
		"06 Глава 6. Бури и\u00a0конфликты":                                                            "06 Глава 6. Бури и конфликты",
		"06 Глава 6. Типы внутренней мотивации":                                                       "06 Глава 6. Типы внутренней мотивации",
		"06 Глава 6. Покончите с\u00a0перегрузкой на\u00a0работе\u00a0—\u00a0 установите границы":      "06 Глава 6. Покончите с перегрузкой на работе —  установите границы",
		"06 Привычка №4. Повышать продуктивность":                                                      "06 Привычка No4. Повышать продуктивность",
		"06. 5 Достижения других людей - пример для вашего развития":                                  "06. 5 Достижения других людей - пример для вашего развития",
		"06. Миф №3. У меня нет способностей к языку":                                                  "06. Миф No3. У меня нет способностей к языку",
		"06. Урок №5. Раскройте свои таланты":                                                          "06. Урок No5. Раскройте свои таланты",
		"06. Глава 5. Как тормозить на\u00a0лифте настроения. Сила любопытства":                        "06. Глава 5. Как тормозить на лифте настроения. Сила любопытства",
		"06. Правило № 3. Мозг каждого человека имеет различную электропроводимость нейронов, часть 2": "06. Правило No 3. Мозг каждого человека имеет различную электропроводимость нейронов, часть 2",
		"06. Принцип №2. Руки убрал быстро!":                                                           "06. Принцип No2. Руки убрал быстро!",
		"07 День 4. Этот чёртов палец":                                                                "07 День 4. Этот чёртов палец",
		"07 Глава 2. Движущийся мир #04":                                                              "07 Глава 2. Движущийся мир #04",
		"07 Глава 6. Великий ковровый путь":                                                          "07 Глава 6. Великий ковровый путь",
		"07 Глава 7. Мы «снимаем сливки» ":                                                             "07 Глава 7. Мы «снимаем сливки»",
		"07 Глава 7. Без вариантов… инжиниринг — для девчонок":                                         "07 Глава 7. Без вариантов... инжиниринг — для девчонок",
		"07 Глава 7. Подъемы и\u00a0спады":                                                             "07 Глава 7. Подъемы и спады",
		"07 Глава 7. Определим цель работы с\u00a0возражениями":                                        "07 Глава 7. Определим цель работы с возражениями",
		"07 Глава 7. Самосознание и\u00a0выбор":                                                        "07 Глава 7. Самосознание и выбор",
		"07 Привычка №5. Оказывать влияние":                                                            "07 Привычка No5. Оказывать влияние",
		"07. Миф №4. Это очень дорого и долго":                                                         "07. Миф No4. Это очень дорого и долго",
		"07. Урок №6. Сопротивление, или Что нас останавливает":                                        "07. Урок No6. Сопротивление, или Что нас останавливает",
		"07. Правило № 4. Мы не обращаем внимания на скучное, часть 1":                                 "07. Правило No 4. Мы не обращаем внимания на скучное, часть 1",
		"07. Принцип №3. Соринка для профессора, бревно для обезьянки":                                 "07. Принцип No3. Соринка для профессора, бревно для обезьянки",
		"08 Глава 2. Движущийся мир #05":                                                              "08 Глава 2. Движущийся мир #05",
		"08 Глава 8 Черный, белый и серый":                                                          "08 Глава 8 Черный, белый и серый",
		"08 Глава 8. Мы любим порядок ":                                                                "08 Глава 8. Мы любим порядок",
		"08 Глава 8. Как спокойно и\u00a0с\u00a0достоинством проститься со страстью":                   "08 Глава 8. Как спокойно и с достоинством проститься со страстью",
		"08 Глава 8. Чем различается общение в\u00a0обиходе и\u00a0в\u00a0бизнесе":                     "08 Глава 8. Чем различается общение в обиходе и в бизнесе",
		"08 Глава 8. Утраты и\u00a0предел возможностей":                                                "08 Глава 8. Утраты и предел возможностей",
		"08 Глава 8. Человек\u00a0— троянский конь":                                                    "08 Глава 8. Человек — троянский конь",
		"08 Часть четвёртая. Дары":                                                                    "08 Часть четвёртая. Дары",
		"08 Привычка №6. Проявлять смелость":                                                           "08 Привычка No6. Проявлять смелость",
		"08. Миф №5. Язык можно выучить, только проживая за границей в течение хотя бы полугода":       "08. Миф No5. Язык можно выучить, только проживая за границей в течение хотя бы полугода",
		"08. Урок №7. Банк идей, или Поиск единомышленников":                                           "08. Урок No7. Банк идей, или Поиск единомышленников",
		"08. Правило № 4. Мы не обращаем внимания на скучное, часть 2":                                 "08. Правило No 4. Мы не обращаем внимания на скучное, часть 2",
		"08. Принцип №4. Снова сегодня":                                                                "08. Принцип No4. Снова сегодня",
		"09 09 Глава 9. Сцепление — твой друг":                                                        "09 09 Глава 9. Сцепление — твой друг",
		"09 Часть 1. Актер. 7. Мыслительный процесс":                                                  "09 Часть 1. Актер. 7. Мыслительный процесс",
		"09 Глава 2. Движущийся мир #06":                                                              "09 Глава 2. Движущийся мир #06",
		"09 Глава 8. Откройте своего партнера":                                                        "09 Глава 8. Откройте своего партнера",
		"09 Глава 9. Как извлечь максимальную пользу из\u00a0отпуска":                                  "09 Глава 9. Как извлечь максимальную пользу из отпуска",
		"09 Глава 9. Какова первичная реакция на\u00a0возражение":                                      "09 Глава 9. Какова первичная реакция на возражение",
		"09 Пробуйте": "09 Пробуйте",
		"09 Заключение. Как жить со\u00a0страстью продуктивно":                                                                      "09 Заключение. Как жить со страстью продуктивно",
		"09. Урок №8. Генеральная репетиция":                                                                                        "09. Урок No8. Генеральная репетиция",
		"09. Правило № 5. Кратковременная память: повторить, чтобы вспомнить, часть 1":                                              "09. Правило No 5. Кратковременная память: повторить, чтобы вспомнить, часть 1",
		"09. Принцип №5. Покажи мне результаты":                                                                                     "09. Принцип No5. Покажи мне результаты",
		"1. Что сложнее, чем найти новые ответы?":                                                                                  "1. Что сложнее, чем найти новые ответы?",
		"1. Продавец — это звучит гордо! ":                                                                                          "1. Продавец — это звучит гордо!",
		"1.13. История успеха № 1. Новые традиции Lockheed Martin, часть 1":                                                         "1.13. История успеха No 1. Новые традиции Lockheed Martin, часть 1",
		"1.14. История успеха № 1. Новые традиции Lockheed Martin, часть 2":                                                         "1.14. История успеха No 1. Новые традиции Lockheed Martin, часть 2",
		"1.19. История успеха № 2. Bank One не только крупнее, но и лучше, часть 1":                                                 "1.19. История успеха No 2. Bank One не только крупнее, но и лучше, часть 1",
		"1.20. История успеха № 2. Bank One не только крупнее, но и лучше, часть 2":                                                 "1.20. История успеха No 2. Bank One не только крупнее, но и лучше, часть 2",
		"1.21. История успеха № 2. Bank One не только крупнее, но и лучше, часть 3":                                                 "1.21. История успеха No 2. Bank One не только крупнее, но и лучше, часть 3",
		"1.31. История успеха № 3. Форт-Уэйн, Индиана, подъем с 0 до 60 проектов бережливого производства в мгновение ока, часть 1": "1.31. История успеха No 3. Форт-Уэйн, Индиана, подъем с 0 до 60 проектов бережливого производства в мгновение ока, часть 1",
		"1.32. История успеха № 3. Форт-Уэйн, Индиана, подъем с 0 до 60 проектов бережливого производства в мгновение ока, часть 2": "1.32. История успеха No 3. Форт-Уэйн, Индиана, подъем с 0 до 60 проектов бережливого производства в мгновение ока, часть 2",
		"1.40. История успеха № 4. Stanford hospital and clinics: на переднем крае революции качества, часть 1":                     "1.40. История успеха No 4. Stanford hospital and clinics: на переднем крае революции качества, часть 1",
		"1.41. История успеха № 4. Stanford hospital and clinics: на переднем крае революции качества, часть 2":                     "1.41. История успеха No 4. Stanford hospital and clinics: на переднем крае революции качества, часть 2",
		"10 День 7. Проверка показателей":                                                                                          "10 День 7. Проверка показателей",
		"10 Глава 10. Типы внешней мотивации":                                                                                      "10 Глава 10. Типы внешней мотивации",
		"10 Глава 10. Пары, у\u00a0которых все получилось":                                                                          "10 Глава 10. Пары, у которых все получилось",
		"10 Глава 10. Социальный мозг в\u00a0жизни":                                                                                 "10 Глава 10. Социальный мозг в жизни",
		"10 Глава 2. Движущийся мир #07":                                                                                           "10 Глава 2. Движущийся мир #07",
		"10 Глава 9. Полезные лайфхаки для начинающих блогеров":                                                                    "10 Глава 9. Полезные лайфхаки для начинающих блогеров",
		"10. Урок №9. Колода памяти и Колода желаний":                                                                               "10. Урок No9. Колода памяти и Колода желаний",
		"10. Глава 3. Сначала кто… затем что, часть 1":                                                                              "10. Глава 3. Сначала кто... затем что, часть 1",
		"10. Делай, как если бы":                                                                                                   "10. Делай, как если бы",
		"10. Правило № 5. Кратковременная память: повторить, чтобы вспомнить, часть 2":                                              "10. Правило No 5. Кратковременная память: повторить, чтобы вспомнить, часть 2",
		"10. Принцип №6. Игры, веселье и все такое":                                                                                 "10. Принцип No6. Игры, веселье и все такое",
		"11 Глава 10. Определите свой курс":                                                                                        "11 Глава 10. Определите свой курс",
		"11 Глава 11. Пять стратегий успешной работы с\u00a0неполной занятостью":                                                    "11 Глава 11. Пять стратегий успешной работы с неполной занятостью",
		"11 Глава 11. Питание ":                                                                   "11 Глава 11. Питание",
		"11 Глава 11. Социальный мозг в\u00a0бизнесе":                                             "11 Глава 11. Социальный мозг в бизнесе",
		"11 Страсть к дизайну — конкурентное преимущество No 1":                                  "11 Страсть к дизайну — конкурентное преимущество No 1",
		"11.  Глава 3. Сначала кто… затем что, часть 2":                                           "11.  Глава 3. Сначала кто... затем что, часть 2",
		"11. Урок №10. Живите полной жизнью":                                                      "11. Урок No10. Живите полной жизнью",
		"11. Правило № 6. Долговременная память: вспомнить, чтобы повторить, часть 1":             "11. Правило No 6. Долговременная память: вспомнить, чтобы повторить, часть 1",
		"12 Глава 11. Об истории и войне":                                                        "12 Глава 11. Об истории и войне",
		"12 Глава 11. Первый день нашей иммиграции":                                             "12 Глава 11. Первый день нашей иммиграции",
		"12 Глава 12 Искусство войны":                                                            "12 Глава 12 Искусство войны",
		"12 Глава 12. Активная и\u00a0пассивная работа с\u00a0возражениями":                       "12 Глава 12. Активная и пассивная работа с возражениями",
		"12 Глава 12. Социальный мозг в\u00a0образовании":                                         "12 Глава 12. Социальный мозг в образовании",
		"12 Глава 12. Оставаться верным себе ":                                                    "12 Глава 12. Оставаться верным себе",
		"12 Глава 12. Сохранение концентрации при работе из\u00a0дома":                            "12 Глава 12. Сохранение концентрации при работе из дома",
		"12 Бесконечный поиск эффективных методов (и еще восемь стратегий добавления ценности)": "12 Бесконечный поиск эффективных методов (и еще восемь стратегий добавления ценности)",
		"12. Глава 11. Как культивировать в\u00a0себе благодарность":                              "12. Глава 11. Как культивировать в себе благодарность",
		"12. Глава 3. Сначала кто… затем что, часть 3":                                            "12. Глава 3. Сначала кто... затем что, часть 3",
		"12. Правило № 6. Долговременная память: вспомнить, чтобы повторить, часть 2":             "12. Правило No 6. Долговременная память: вспомнить, чтобы повторить, часть 2",
		"13 Дни 10–11. Здравствуй, боль":                                                         "13 Дни 10–11. Здравствуй, боль",
		"13 Глава 12. Портрет двух отношений":                                                    "13 Глава 12. Портрет двух отношений",
		"13 Глава 13. Типы личной мотивации":                                                     "13 Глава 13. Типы личной мотивации",
		"13.  Глава 3. Сначала кто… затем что, часть 4":                                           "13.  Глава 3. Сначала кто... затем что, часть 4",
		"13. Глава 12. Уважение к\u00a0чужой реальности":                                          "13. Глава 12. Уважение к чужой реальности",
		"13. Правило № 7. Хороший сон — хорошее мышление, часть 1":                                "13. Правило No 7. Хороший сон — хорошее мышление, часть 1",
		"14 Глава 14. Игнорирование друзей из-за работы не\u00a0поможет вашей карьере":            "14 Глава 14. Игнорирование друзей из-за работы не поможет вашей карьере",
		"14 Часть 2. Писательский настрой. Смотрим по сторонам":                                 "14 Часть 2. Писательский настрой. Смотрим по сторонам",
		"14 Курсовая устойчивость":                                                               "14 Курсовая устойчивость",
		"14 Линейные менеджеры — самый недооцененный актив":                                    "14 Линейные менеджеры — самый недооцененный актив",
		"14. Глава 13. Развивайте преданность и\u00a0оптимизм":                                    "14. Глава 13. Развивайте преданность и оптимизм",
		"14. Правило № 7. Хороший сон — хорошее мышление, часть 2":                                "14. Правило No 7. Хороший сон — хорошее мышление, часть 2",
		"15 Глава 14. Женский фактор":                                                            "15 Глава 14. Женский фактор",
		"15 Обращение за профессиональной помощью":                                               "15 Обращение за профессиональной помощью",
		"15 Моральный аспект":                                                                    "15 Моральный аспект",
		"15. Глава 14. Что делать, если настроение на\u00a0нуле":                                  "15. Глава 14. Что делать, если настроение на нуле",
		"15. Правило № 8. Стресс негативно влияет на способность мозга учиться, часть 1":          "15. Правило No 8. Стресс негативно влияет на способность мозга учиться, часть 1",
		"16 День 14. Проблемы с математикой":                                                     "16 День 14. Проблемы с математикой",
		"16 Глава 15. Сахалинский период":                                                        "16 Глава 15. Сахалинский период",
		"16 Глава 16. Пять причин, почему возражение\u00a0—\u00a0 это зд'орово":                   "16 Глава 16. Пять причин, почему возражение —  это зд'орово",
		"16 Глава 16. Типы социальной мотивации":                                                 "16 Глава 16. Типы социальной мотивации",
		"16 Часть 2. Сценические упражнения. 13. Текущий момент":                                 "16 Часть 2. Сценические упражнения. 13. Текущий момент",
		"16. Глава 15. Отношения и\u00a0лифт настроения":                                          "16. Глава 15. Отношения и лифт настроения",
		"16. Правило № 8. Стресс негативно влияет на способность мозга учиться, часть 2":          "16. Правило No 8. Стресс негативно влияет на способность мозга учиться, часть 2",
		"17 День 15. Отжимания - наше всё":                                                       "17 День 15. Отжимания - наше всё",
		"17 Отстой FM": "17 Отстой FM",
		"17. Правило № 9. Сенсорная интеграция: задействуйте больше чувств, часть 1": "17. Правило No 9. Сенсорная интеграция: задействуйте больше чувств, часть 1",
		"17. Отличайся!":                   "17. Отличайся!",
		"18 Глава 18. Сочетание мотиваций": "18 Глава 18. Сочетание мотиваций",
		"18. Правило № 9. Сенсорная интеграция: задействуйте больше чувств, часть 2":                  "18. Правило No 9. Сенсорная интеграция: задействуйте больше чувств, часть 2",
		"19 Глава 19. Шесть способов включить заботу о\u00a0себе в рабочий день":                      "19 Глава 19. Шесть способов включить заботу о себе в рабочий день",
		"19 Глава 19. Вопросы, которые следует задать, прежде чем начинать любой творческий проект": "19 Глава 19. Вопросы, которые следует задать, прежде чем начинать любой творческий проект",
		"19. Мы продаем будущий образ жизни Клиента. И отвечаем за это":                              "19. Мы продаем будущий образ жизни Клиента. И отвечаем за это",
		"19. Правило № 10. Зрение важнее остальных сенсорных органов, часть 1":                        "19. Правило No 10. Зрение важнее остальных сенсорных органов, часть 1",
		"2. Сопли — это не геройство":                                                                "2. Сопли — это не геройство",
		"20 Глава 20. Пять моделей развития карьеры для творческого человека":                        "20 Глава 20. Пять моделей развития карьеры для творческого человека",
		"20 Глава 20. Разница между трудоголизмом и\u00a0работой сверхурочно":                         "20 Глава 20. Разница между трудоголизмом и работой сверхурочно",
		"20 Техника 6. Сила «хм…»": "20 Техника 6. Сила «хм...»",
		"20. Правило № 10. Зрение важнее остальных сенсорных органов, часть 2":                     "20. Правило No 10. Зрение важнее остальных сенсорных органов, часть 2",
		"21 Глава 21 Склейте меня обратно":                                                        "21 Глава 21 Склейте меня обратно",
		"21 Глава 21. Как забыть о\u00a0работе, когда вы не\u00a0работаете":                        "21 Глава 21. Как забыть о работе, когда вы не работаете",
		"21. Правило № 11. Гендер: мозг мужчины и женщины различен, часть 1":                       "21. Правило No 11. Гендер: мозг мужчины и женщины различен, часть 1",
		"22 Глава 22. Модели управления временем ":                                                 "22 Глава 22. Модели управления временем",
		"22 Техника 8. От транзакции к трансформации ":                                             "22 Техника 8. От транзакции к трансформации",
		"22. Правило № 11. Гендер: мозг мужчины и женщины различен, часть 2":                       "22. Правило No 11. Гендер: мозг мужчины и женщины различен, часть 2",
		"23 День 21. Начинай по секундной стрелке":                                               "23 День 21. Начинай по секундной стрелке",
		"23 Часть 2. Сценические упражнения. 20. Действия героя":                                  "23 Часть 2. Сценические упражнения. 20. Действия героя",
		"23 Глава 23 Ползком к финишной черте":                                                    "23 Глава 23 Ползком к финишной черте",
		"23 Глава 23. Начните контролировать свою поездку на\u00a0работу":                          "23 Глава 23. Начните контролировать свою поездку на работу",
		"23. Правило № 12. Исследование: по своей природе мы великие первооткрыватели, часть 1":    "23. Правило No 12. Исследование: по своей природе мы великие первооткрыватели, часть 1",
		"24 Глава 24. Еще один активный инструмент-«Обязательство» ":                               "24 Глава 24. Еще один активный инструмент-«Обязательство»",
		"24 Творческий кризис":                                                                    "24 Творческий кризис",
		"24. Всегда управляй сроками!":                                                            "24. Всегда управляй сроками!",
		"24. Правило № 12. Исследование: по своей природе мы великие первооткрыватели, часть 2":    "24. Правило No 12. Исследование: по своей природе мы великие первооткрыватели, часть 2",
		"25 Глава 25. «Закрытие»\u00a0—\u00a0 инструмент продаж, переговоров и любой коммуникации": "25 Глава 25. «Закрытие» —  инструмент продаж, переговоров и любой коммуникации",
		"25 Часть 3. Пьеса и роль. 21. Первый контакт с пьесой":                                  "25 Часть 3. Пьеса и роль. 21. Первый контакт с пьесой",
		"25. Пока есть хоть один шанс, бейся до конца!":                                           "25. Пока есть хоть один шанс, бейся до конца!",
		"26 День 25. По самое «не балуйся»":                                                       "26 День 25. По самое «не балуйся»",
		"26 Глава 23. Пусть говорят - фразы, которые вы наверняка услышите ":                       "26 Глава 23. Пусть говорят - фразы, которые вы наверняка услышите",
		"26. Люди покупают у людей":                                                               "26. Люди покупают у людей",
		"27 Глава 27. Как общаться с\u00a0компетентным клиентом":                                   "27 Глава 27. Как общаться с компетентным клиентом",
		"28 С днём рождения!":                                                                     "28 С днём рождения!",
		"28 День 27. Тысяча отжиманий":                                                            "28 День 27. Тысяча отжиманий",
		"28. Продажи делаются по одной":                                                           "28. Продажи делаются по одной",
		"29 Часть 5. Последний урок":                                                              "29 Часть 5. Последний урок",
		"3 - Глава вторая. «Как вы себя чувствуете_», или Осознание и понимание своих эмоций":     "3 - Глава вторая. «Как вы себя чувствуете_», или Осознание и понимание своих эмоций",
		"3. Каждый из нас продавец радости!":                                                      "3. Каждый из нас продавец радости!",
		"30 Глава 30. Как попасть аргументом в\u00a0цель":                                          "30 Глава 30. Как попасть аргументом в цель",
		"31 Если небо упадёт":                                                                     "31 Если небо упадёт",
		"31 Часть 3. Пьеса и роль. 27. Действие":                                                  "31 Часть 3. Пьеса и роль. 27. Действие",
		"32. День 31. Грустный деньmp3":                                                           "32. День 31. Грустный деньmp3",
		"34 Сарай": "34 Сарай",
		"34. Верить в светлое завтра, а действовать от наихудшего для себя варианта": "34. Верить в светлое завтра, а действовать от наихудшего для себя варианта",
		"35 Глава 35. В\u00a0чем помогает собственный опыт":                           "35 Глава 35. В чем помогает собственный опыт",
		"38 Глава 38. Комплимент как мощнейший инструмент ":                           "38 Глава 38. Комплимент как мощнейший инструмент",
		"38. Не обижайте лучших продавцов":                                           "38. Не обижайте лучших продавцов",
		"39 Глава 39. Рефрейминг и\u00a0общность\u00a0—\u00a0 как их использовать":    "39 Глава 39. Рефрейминг и общность —  как их использовать",
		"4 - Глава третья. Осознание и понимание эмоций других ":                     "4 - Глава третья. Осознание и понимание эмоций других",
		"40 1504. ФЛОРЕНЦИЯ Микеланджело. Май 1504 года. Флоренция":                  "40 1504. ФЛОРЕНЦИЯ Микеланджело. Май 1504 года. Флоренция",
		"40 Прощай": "40 Прощай",
		"40. Не втягивайте Клиентов в свои внутренние проблемы":                             "40. Не втягивайте Клиентов в свои внутренние проблемы",
		"41 Глава 41. Как отрабатывать «В\u00a0другом месте дешевле»":                        "41 Глава 41. Как отрабатывать «В другом месте дешевле»",
		"42 Глава 42. Простая отработка отговорок и\u00a0отказа":                             "42 Глава 42. Простая отработка отговорок и отказа",
		"43. Невыставленный счет не может быть оплачен ":                                    "43. Невыставленный счет не может быть оплачен",
		"44 Глава завершающая. Аргументация\u00a0—\u00a0 это не\u00a0торг":                   "44 Глава завершающая. Аргументация —  это не торг",
		"5 - Глава четвертая. «Учитесь властвовать собой», или Управление своими эмоциями ": "5 - Глава четвертая. «Учитесь властвовать собой», или Управление своими эмоциями",
		"5. Хочешь понравиться — сиди и слушай":                                             "5. Хочешь понравиться — сиди и слушай",
		"58. Глава 11. Обретение видения, часть 4 ":                                          "58. Глава 11. Обретение видения, часть 4",
		"6. Слишком хороший контакт с Клиентом уменьшает вероятность продажи":               "6. Слишком хороший контакт с Клиентом уменьшает вероятность продажи",
		"60. Глава 11. Обретение видения, часть 6 ":                                          "60. Глава 11. Обретение видения, часть 6",
		"7. Основной мотив наших людей — желание быть героем":                              "7. Основной мотив наших людей — желание быть героем",
		"8. Подготовка — лучший друг продавана":                                             "8. Подготовка — лучший друг продавана",
		"9. Дайте попробовать. Это всегда работает!":                                        "9. Дайте попробовать. Это всегда работает!",
		"В этом году я…":                                                   "В этом году я...",
		"Папа — морской конек":                                            "Папа — морской конек",
		"Глава 1. Понимаем сумасшедших ":                                   "Глава 1. Понимаем сумасшедших",
		"Глава 10. На одной волне":                                        "Глава 10. На одной волне",
		"Глава 10. Не-делай-это-сам":                                      "Глава 10. Не-делай-это-сам",
		"Глава 10. Перегруппировка —\u00a0начинаем новую игру":             "Глава 10. Перегруппировка — начинаем новую игру",
		"Глава 11. Слепящий свет, часть I":                                "Глава 11. Слепящий свет, часть I",
		"Глава 12. Слепящий свет, часть II":                               "Глава 12. Слепящий свет, часть II",
		"Глава 13. 1919-й вернулся!":                                      "Глава 13. 1919-й вернулся!",
		"Глава 19. Решительный отказ и вежливый отказ":                   "Глава 19. Решительный отказ и вежливый отказ",
		"Глава 2. Окружение, создающее победителей":                       "Глава 2. Окружение, создающее победителей",
		"Глава 2. Распознаём механизм безумия":                            "Глава 2. Распознаём механизм безумия",
		"Глава 3. Как определить образ действия иррационального человека": "Глава 3. Как определить образ действия иррационального человека",
		"Глава 32. Что делать, если ваш близкий хочет покончить с собой": "Глава 32. Что делать, если ваш близкий хочет покончить с собой",
		"Глава 4. Дочь Мари Кэтрин Мэри Эллен О’Брайен Фрид":              "Глава 4. Дочь Мари Кэтрин Мэри Эллен О’Брайен Фрид",
		"Глава 5. Распознаём внутреннего сумасшедшего":                    "Глава 5. Распознаём внутреннего сумасшедшего",
		"Глава 6. Я слишком многого хочу —\u00a0и разбрасываюсь":           "Глава 6. Я слишком многого хочу — и разбрасываюсь",
		"Глава 6. Мозговой штурм":                                         "Глава 6. Мозговой штурм",
		"Глава 7. Это война, вездесущая и бесконечная":                    "Глава 7. Это война, вездесущая и бесконечная",
		"Глава 7. Жизнь —\u00a0это удаленный сотрудник":                    "Глава 7. Жизнь — это удаленный сотрудник",
		"Шахматы с енотом. Рабочая тетрадь № 1":                            "Шахматы с енотом. Рабочая тетрадь No 1",
		"Шахматы с енотом. Рабочая тетрадь № 2":                            "Шахматы с енотом. Рабочая тетрадь No 2",
		"Шахматы с енотом. Рабочая тетрадь № 3":                            "Шахматы с енотом. Рабочая тетрадь No 3",
		"Тайные соседи":                                                   "Тайные соседи",
		"Частый гребешок":                                                 "Частый гребешок",
		"Папагай и Пашагай":                                              "Папагай и Пашагай",
		"Львиный рык":                                                     "Львиный рык",
		"Воздушный змей":                                                 "Воздушный змей",
		"Спокойной ночи!":                                                "Спокойной ночи!",
		"Вкуснейский павиан и яичница":                                   "Вкуснейский павиан и яичница",
		"Эмоциональный интеллект для\u00a0менеджеров проектов":             "Эмоциональный интеллект для менеджеров проектов",
		"Эмоциональный интеллект":                                         "Эмоциональный интеллект",
	}
	for arg, want := range tests {
		t.Run(arg, func(t *testing.T) {
			require.Equal(t, want, String(arg))
		})
	}
}