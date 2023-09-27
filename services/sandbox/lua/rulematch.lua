local utils = require 'utils'
local w32 = require 'w32api'
local sysevent = require 'sysevents'
local rules = require 'rules.init'

---@class Context Context structure for this run.
---@field reg_key_handles table<string,string> Track reg key handle -> reg key path.
---@field file_handles table<string,string> Track file handle -> reg file path.
---@field network_handles table<string,string> Track netwotk handle -> network address.
---@field events table<integer,string> Summarized system events.
---@field seen_events  table<string,boolean> Seen events.
---@field created_files string  The operation requested over the above `path` field.
local Context = {}

function Context:new(o)
    o = o or {}
    setmetatable(o, self)
    self.__index = self
    o.reg_key_handles = {}
    o.file_handles = {}
    o.network_handles = {}
    o.events = {}
    o.seen_events = {}
    return o
end

local M = {}

--- Evaluate an API trace against a set of rules.
---@param w32apis table<integer, Win32API> Table containing list of API objects.
function M:eval(w32apis, artifacts)

    local ctx = Context:new()
    ctx.reg_key_handles = utils.table_copy(w32.pre_reg_key_handles)

    ---@type table<string, table<integer, Win32API>>
    local APINamesTable = {}

    for _, w32api in ipairs(w32apis) do

        -- Index all APIs on their names for quick lookups.
        if APINamesTable[w32api.name] == nil then
            APINamesTable[w32api.name] = {}
            table.insert(APINamesTable[w32api.name], w32api)
        else
            table.insert(APINamesTable[w32api.name], w32api)
        end

        -- Extract system events.
        local event = sysevent:summarize(ctx, w32api)
        if event ~= nil then
            print(event)
        end
    end

    -- Run each rule against the API trace.
    for _, rule in pairs(rules) do
        local res = rule:Match(ctx)
    end

    ctx = nil

    return nil
end

return M
