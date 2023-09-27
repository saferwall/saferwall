local ANNOTATION = {IN = 'in', OUT = 'out', RESERVED = 'reserved'}

---@class Win32APIParameter
---@field annotation string SAL annotation
---@field name string Argument name
---@field value any Argument value
local Win32APIParameter = {}

function Win32APIParameter:new(o, anno, name, val)
    o = o or {}
    setmetatable(o, self)
    self.__index = self
    o.annotation = anno
    o.name = name
    o.value = val
    return o
end

---@class Win32API
---@field timestamp number Timestamp when the API was executed
---@field process_id string Process identifier
---@field thread_id string Thread identifier
---@field name string API name
---@field parameters table<integer, Win32APIParameter> List of Parameters
---@field return_value string Return value
local Win32API = {}

function Win32API:new(o, ts, pid, tid, name, params, ret_value)
    o = o or {}
    setmetatable(o, self)
    self.__index = self
    o.timestamp = ts
    o.process_id = pid
    o.thread_id = tid
    o.name = name
    o.return_value = ret_value
    o.parameters = {}
    if type(params) == 'userdata' then
        params = {}
    end
    for _, param in ipairs(params or {}) do
        table.insert(o.parameters, Win32APIParameter:new(nil, param.anno, param.name, param.val))
    end
    return o
end

function Win32API:param_value_from_name(param_name)
    for _, param in ipairs(self.parameters) do
        if param.name == param_name then
            if param.annotation == ANNOTATION.IN or param.annotation == ANNOTATION.OUT or param.annotation ==
                ANNOTATION.RESERVED then
                return param.value
            end
        end
    end

    return nil
end

return Win32API
