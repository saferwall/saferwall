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

--- Filter a table using a predicate function.
function M.table_filter(func, t)
    local rettab = {}
    for _, entry in pairs(t) do
        if func(entry) then
            table.insert(rettab, entry)
        end
    end
    return rettab
end

--- Create a flat list of all files in a directory
-- @param directory - The directory to scan (default value = './')
-- @param recursive - Whether or not to scan subdirectories recursively (default value = true)
-- @param extensions - List of extensions to collect, if blank all will be collected
function M.scan_dir(directory, recursive, extensions)
    directory = directory or ''
    recursive = recursive or true

    local currentDirectory = directory
    local fileList = {}
    local command = 'ls ' .. currentDirectory .. ' -p'

    -- if string.sub(directory, -1) ~= '/' then directory = directory .. '/' end
    if recursive then
        command = command .. 'R'
    end

    for fileName in io.popen(command):lines() do
        if string.sub(fileName, -1) == '/' then
            -- Directory, don't do anything
        elseif string.sub(fileName, -1) == ':' then
            currentDirectory = string.sub(fileName, 1, -2)
            -- if currentDirectory ~= directory then
            currentDirectory = currentDirectory .. '/'
            -- end
        elseif string.len(fileName) == 0 then
            -- Blank line
            currentDirectory = directory
            -- elseif string.find(fileName,"%.lua$") then
            -- File is a .lua file
        else
            if type(extensions) == 'table' then
                for _, extension in ipairs(extensions) do
                    if string.find(fileName, '%.' .. extension .. '$') then
                        table.insert(fileList, currentDirectory .. fileName)
                    end
                end
            else
                table.insert(fileList, currentDirectory .. fileName)
            end
        end
    end

    return fileList
end

function M.path_base_name(file)
    return file:match('^.+/(.+)$')
end

function M.split(pString, pPattern)
    local Table = {} -- NOTE: use {n = 0} in Lua-5.0
    local fpat = '(.-)' .. pPattern
    local last_end = 1
    local s, e, cap = pString:find(fpat, 1)
    while s do
        if s ~= 1 or cap ~= '' then
            table.insert(Table, cap)
        end
        last_end = e + 1
        s, e, cap = pString:find(fpat, last_end)
    end
    if last_end <= #pString then
        cap = pString:sub(last_end)
        table.insert(Table, cap)
    end
    return Table
end

return M
