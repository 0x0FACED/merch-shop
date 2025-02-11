-- Создаем схему для группировки всех таблиц сервиса
CREATE SCHEMA IF NOT EXISTS shop;

-- Таблицы users просто хранит имя пользователя (уникальное) и хэш пароля
-- Хэш создаем через bcrypt (самый просто вариант изначально), далее возможен переход на argon2, например
CREATE TABLE IF NOT EXISTS shop.users (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL, -- Оставил размер 255, потому что нет ограничения в спецификации
    password_hash TEXT NOT NULL, -- делаем TEXT, потому что алгортим шифрования может поменяться, как и длина хэша
    created_at TIMESTAMP DEFAULT NOW()
);

-- Таблица для хранения кошельков юзеров. Создается автоматически при создании юзера
-- PRIMARY KEY ссылается на приватный ключ users, при удалении юзера автоматически удалится и его кошелек
-- Баланс по дефолту 0, CHECK для того, чтобы в минус нельзя было уйти
CREATE TABLE IF NOT EXISTS shop.wallets (
    user_id INTEGER PRIMARY KEY REFERENCES shop.users(id) ON DELETE CASCADE,
    balance INTEGER NOT NULL DEFAULT 1000 CHECK (balance >= 0)
);

-- Таблица для хранения всех транзакций
-- При удалении юзера юзера транзакция не удалится, а user_id станет NULL (юзера не существует), но транзу надо сохранить.
-- amount, соответственно, сумма. Не может быть меньше 0 (для этого есть CHECK)
CREATE TABLE IF NOT EXISTS shop.transactions (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    from_user_id INTEGER REFERENCES shop.users(id) ON DELETE SET NULL,
    to_user_id INTEGER REFERENCES shop.users(id) ON DELETE SET NULL,
    amount INTEGER NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Таблица для хранения всех предметов, доступных в магазине
-- Подразумевается расширение магазина в дальнейшем через какую-нибудь админ панель, к которой будет иметь доступ, удивительно, админы
-- Они смогут добавлять предметы.
-- У предмета есть название и цена. Цена ТОЛЬКО INT (ну для чего копейки нужны, for real?)
CREATE TABLE IF NOT EXISTS shop.items (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0)
);

-- Инвентарь представляет из себя просто запись о предмете, его количестве и владельце
-- Например, пользователь купил футболку "I love Bookks" в количестве 10шт, тогда
-- создастся запись, где user_id будет наш юзер, item_id - id этой футболки, а 
-- qunatity - количество этих футболок
-- Первичный ключ составной (зависит от юзера и предмета напрмую) 
CREATE TABLE IF NOT EXISTS shop.inventory (
    user_id INTEGER REFERENCES shop.users(id) ON DELETE CASCADE,
    item_id INTEGER REFERENCES shop.items(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    PRIMARY KEY (user_id, item_id)
);

-- просто индексы для ускорения поиска
CREATE INDEX IF NOT EXISTS idx_users_username ON shop.users(username);

CREATE INDEX IF NOT EXISTS idx_wallets_balance ON shop.wallets(balance);

CREATE INDEX IF NOT EXISTS idx_items_name ON shop.items(name);

CREATE INDEX IF NOT EXISTS idx_transactions_from_user ON shop.transactions(from_user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_to_user ON shop.transactions(to_user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_from_to ON shop.transactions(from_user_id, to_user_id);

CREATE INDEX IF NOT EXISTS idx_inventory_user ON shop.inventory(user_id);
CREATE INDEX IF NOT EXISTS idx_inventory_item ON shop.inventory(item_id);
CREATE INDEX IF NOT EXISTS idx_inventory_user_item_quantity ON shop.inventory(user_id, item_id, quantity);

INSERT INTO shop.items (name, price) VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);