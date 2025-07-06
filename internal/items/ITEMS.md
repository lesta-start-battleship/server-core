Начало игры (Загружаем скрипты предметов) - `items, err := GetAllItems()`

* []Item - список доступных предметов

Перед использованием предмета (Загружаем инвентарь) - `inv, err := GetNumberItems(token)`

* token string - JWT токен пользователя
* map[ItemID]int - мапа предметов в инвентаре

Применяем предмет - `UseItem(itemID, state, itemsList, params, userJWT)`

* itemID ItemID - айди предмета который нужно использовать
* state *game.States - состояние игры
* itemsList map[ItemID]*Item - предметы
* params map[string]interface{} - дополнительные параметры для скрипта
* userJWT string - JWT токен пользователя

