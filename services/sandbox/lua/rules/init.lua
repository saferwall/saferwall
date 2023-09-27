local utils = require 'utils'

-- Require all other `.lua` files in the same directory
local info = debug.getinfo(1, 'S')
local module_directory = string.match(info.source, '^@(.*)/')
local module_filename = string.match(info.source, '/([^/]*)$')

-- Apparently the name of this module is given as an argument when it is
-- required, and apparently we get that argument with three dots.
local module_name = ...
module_name = module_name:gsub('.init', '')

-- Walk recursively all files within the required base module and
-- exclude all non lua files + this init.lua module itself.
local module_paths = utils.table_filter(function(filename)
    local is_lua_module = string.match(filename, '[.]lua$')
    local is_this_file = utils.path_base_name(filename) == module_filename
    return is_lua_module and not is_this_file
end, utils.scan_dir(module_directory, true))

-- Load each module.
local ret = {}
for _, module_path in ipairs(module_paths) do
    local path_elements = utils.split(module_path, package.config:sub(1, 1))

    -- Convert the module file path to a dot format for require().
    local submodule_name = ''
    local found = false
    for _, element in ipairs(path_elements) do
        if element == module_name then
            found = true
        end
        if found then
            submodule_name = submodule_name .. '.' .. element
        end
    end

    submodule_name = submodule_name:sub(2):gsub('.lua', '')
    print('loading ' .. submodule_name)
    local success, mod = pcall(require, submodule_name)

    if success then
        ret[submodule_name] = mod
    else
        print('error loading module \'%s\': %s', submodule_name, mod)
    end
end

return ret
