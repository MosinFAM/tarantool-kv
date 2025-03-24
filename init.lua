box.cfg{
    listen = 3301
}

-- Создание пространства и индекса
if not box.space.kv then
    box.schema.space.create('kv', {
        format = {
            {name = 'key', type = 'string'},
            {name = 'value', type = 'string'}
        }
    })
    box.space.kv:create_index('primary', {type = 'hash', parts = {'key'}})
end

-- Функция вставки
function insert_kv(key, value)
    if box.space.kv:get(key) then
        error("key already exists")
    end
    return box.space.kv:insert{key, value}
end

-- Функция получения значения
function get_kv(key)
    local result = box.space.kv:get(key)
    if result then
        return result
    else
        return nil, "key not found"
    end
end

-- Функция удаления
function delete_kv(key)
    return box.space.kv:delete(key)
end

-- Функция обновления
function update_kv(key, value)
    if not box.space.kv:get(key) then
        return nil, "key not found"
    end
    return box.space.kv:put{key, value}
end

-- Регистрация функций

box.schema.func.create('insert_kv')
box.schema.func.create('get_kv')
box.schema.func.create('update_kv')
box.schema.func.create('delete_kv')

-- Права гостю на выполнение этих функций

box.schema.user.grant('guest', 'execute', 'function', 'insert_kv')
box.schema.user.grant('guest', 'execute', 'function', 'get_kv')
box.schema.user.grant('guest', 'execute', 'function', 'update_kv')
box.schema.user.grant('guest', 'execute', 'function', 'delete_kv')

-- Права гостю на чтение и запись в space.kv

box.schema.user.grant('guest', 'read,write', 'space', 'kv')

print("Tarantool KV storage initialized")