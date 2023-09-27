local M = {}

function M.table_merge(...)
    local tables_to_merge = {...}
    assert(#tables_to_merge > 1, 'There should be at least two tables to merge them')

    for k, t in ipairs(tables_to_merge) do
        assert(type(t) == 'table', string.format('Expected a table as function parameter %d', k))
    end

    local result = tables_to_merge[1]

    for i = 2, #tables_to_merge do
        local from = tables_to_merge[i]
        for k, v in pairs(from) do
            if type(k) == 'number' then
                table.insert(result, v)
            elseif type(k) == 'string' then
                if type(v) == 'table' then
                    result[k] = result[k] or {}
                    result[k] = M.table_merge(result[k], v)
                else
                    result[k] = v
                end
            end
        end
    end

    return result
end

function M.table_contains(tbl, x)
    local found = false
    for _, v in pairs(tbl) do
        if v == x then
            found = true
        end
    end
    return found
end

function M.table_copy(obj, seen)
    -- Handle non-tables and previously-seen tables.
    if type(obj) ~= 'table' then
        return obj
    end
    if seen and seen[obj] then
        return seen[obj]
    end

    -- New table; mark it as seen an copy recursively.
    local s = seen or {}
    local res = {}
    s[obj] = res
    for k, v in next, obj do
        res[M.table_copy(k, s)] = M.table_copy(v, s)
    end
    return setmetatable(res, getmetatable(obj))
end

return M
