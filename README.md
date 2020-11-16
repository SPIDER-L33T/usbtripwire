# usbtripwire
Small USB-based tripwire for your server/PC

## EN
### What is this for?
1. "Weapon of the last chance". If the bandits torture you and ask for secret data from your computer, tell them to insert a "secret flash disk" into the computer. As soon as they insert the flash into the computer, the script that you have prepared in advance will be executed (for example: remove the hidden directory from the computer) and, if you specify the token in advance, a message will be sent via Telegram to your partner.
2. "Ordinary tripwire". Someone in your absence decided to connect some USB device to a computer or server (for example: to copy data). As soon as he does this, UsbTripwire will send you a notification in Telegram.

### Install

Run command as root:
```bash
wget https://github.com/SPIDER-L33T/usbtripwire/releases/download/v1.0/setup
chmod +x setup
./setup
```
And follow the instructions on the screen.

If you do not specify a secret device during installation, then the program will work as a tripwire that reacts to any installed USB device

After installing run:
```
systemctl restart usbtripwire
```
### Thanks

Thank you for your moral support [Defcon community](https://defcon.org/) and [DC7499](https://defcon.su/)



## RUS
### Для чего это?
1. "Оружие последнего шанса". Если бандиты пытают Вас и спрашивают секретные данные с Вашего компьютера, скажите им, что нужно вставить "секретный флэш-диск" в компьютер. Как только они вставят флэш в компьютер, выполнится скрипт, который Вы заранее подготовите  (например: удалит с компьютера скрытую директорию) и, если заранее укажете токен, отправится сообщение через Телеграм Вашему партнеру.
2. "Обыкновенная растяжка". Кто-то в Ваше отсутствие решил подключить к компьютеру или серверу какое-нибудь УСБ-устройство (например: чтобы скопировать данные). Как только он это сделает, UsbTripWire отправит Вам уведомление в Телеграм.

### Установка

Выполните команды от пользователя root:
```bash
wget https://github.com/SPIDER-L33T/usbtripwire/releases/download/v1.0/setup
chmod +x setup
./setup
```
И следуйте инструкциям на экране.

Если при установке вы не укажете секретное устройство, то программа будет работать как растяжка, реагирующая на любое установленное USB-устройство.

После установки выполните:
```
systemctl restart usbtripwire
```
### Конфигурация

Файл usbtripwire.conf содержит следующие настройки:
```
telegram_apikey=""
```
Ключ API для бота Telegram
```
telegram_alarmtext="Something plug your device into PC"
```
Сообщение, которое бот отправит через Telegram
```
telegram_users=""
```
список chatid (пользователей или каналов), куда нужно отправлять сообщение. В качестве разделителя используется ";" (точка запятой)
```
devlist=""
```
Список USB-устройтв, на появление которых система отреагирует. В качестве разделителя используется ";" (точка запятой).
Если список пуст, то реагировать будет на любое установленное устройство.
```
cmd="date >> /root/log.txt"
```
Команда, которую нужно выполнить при срабатывании. По-умолчанию (для примера) скидывает дату в тектовый файл.
### Благодарность

Спасибо за моральную поддержку [Defcon community](https://defcon.org/) и [DC7499](https://defcon.su/)
