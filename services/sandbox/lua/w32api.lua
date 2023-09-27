local utils = require 'utils'

local M = {
    -- Registry APIs.
    REG_CREATE_APIS = {'RegCreateKeyA', 'RegCreateKeyW', 'RegCreateKeyExA', 'RegCreateKeyExW'},
    REG_OPEN_APIS = {'RegOpenKeyA', 'RegOpenKeyW', 'RegOpenKeyExA', 'RegOpenKeyExW'},
    REG_SET_VALUE_APIS = {'RegSetValueA', 'RegSetValueW', 'RegSetValueExA', 'RegSetValueExW'},
    REG_SET_KEY_VALUE_APIS = {'RegSetKeyValueA', 'RegSetKeyValueW'},
    REG_DELETE_APIS = {
        'RegDeleteKeyA', 'RegDeleteKeyW', 'RegDeleteKeyExA', 'RegDeleteKeyExW', 'RegDeleteValueA', 'RegDeleteValueW'
    },

    -- File APIs.
    FILE_CREATE_FILE_APIS = {
        'CreateFileA', 'CreateFileW', 'CreateFile2', 'CreateFileTransactedA', 'CreateFileTransactedW'
    },
    FILE_CREATE_DIR_APIS = {
        'CreateDirectory', 'CreateDirectoryExA', 'CreateDirectoryExW', 'CreateDirectoryTransactedA',
        'CreateDirectoryTransactedW'
    },
    FILE_OPEN_APIS = {'OpenFile'},
    FILE_READ_APIS = {'ReadFile', 'ReadFileEx'},
    FILE_WRITE_APIS = {'WriteFile', 'WriteFileEx'},
    File_COPY_APIS = {'CopyFileA', 'CopyFileW', 'CopyFileExA', 'CopyFileExW'},
    File_MOVE_APIS = {'MoveFileA', 'MoveFileW', 'MoveFileWithProgressA', 'MoveFileWithProgressW'},
    FILE_DELETE_APIS = {'DeleteFileA', 'DeleteFileW'},

    --- Network APIs.
    NET_WIN_HTTP_APIS = {'WinHttpConnect'},
    NET_WIN_INET_APIS = {'InternetConnectA', 'InternetConnectW'},
    NET_WINSOCK_APIS = {'getaddrinfo', 'GetAddrInfoW', 'GetAddrInfoExA', 'GetAddrInfoExW'},

    -- Default port numbers for WinINet.
    DEFAULT_PORT = {INVALID = 0, FTP = 21, GOPHER = 70, HTTP = 80, HTTPS = 443, SOCKS = 1080},

    REG_SET_APIS = {},
    REGISTRY_APIS = {},
    FILE_CREATE_APIS = {},
    FILE_APIS = {},
    NETWORK_APIS = {},

    --- Predefined registry key handles.
    pre_reg_key_handles = {},

    --- Handle types constants.
    HANDLE_TYPE = {REGISTRY = 1, FILE = 2, NETWORK = 3},

    --- SAL annotations.
    ANNOTATION = {IN = 'in', OUT = 'out', RESERVED = 'reserved'}

}

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
            if param.annotation == M.ANNOTATION.IN or param.annotation == M.ANNOTATION.OUT or param.annotation ==
                M.ANNOTATION.RESERVED then
                return param.value
            end
        end
    end

    return nil
end

function M:handle_to_object_name(ctx, handle, type)
    if type == self.HANDLE_TYPE.FILE then
        return ctx.file_handles[handle]
    elseif type == self.HANDLE_TYPE.REGISTRY then
        return ctx.reg_key_handles[handle]
    elseif type == self.HANDLE_TYPE.NETWORK then
        return ctx.network_handles
    end
end

function M:init()

    utils.table_merge(M.REG_SET_APIS, M.REG_SET_VALUE_APIS, M.REG_SET_KEY_VALUE_APIS)
    utils.table_merge(M.REGISTRY_APIS, M.REG_CREATE_APIS, M.REG_OPEN_APIS, M.REG_SET_APIS, M.REG_DELETE_APIS)
    utils.table_merge(M.FILE_CREATE_APIS, M.FILE_CREATE_FILE_APIS, M.FILE_CREATE_DIR_APIS)
    utils.table_merge(M.FILE_APIS, M.FILE_CREATE_APIS, M.FILE_OPEN_APIS, M.FILE_READ_APIS, M.FILE_WRITE_APIS,
                      M.File_COPY_APIS, M.File_MOVE_APIS, M.FILE_DELETE_APIS)
    utils.table_merge(M.NETWORK_APIS, M.NET_WIN_INET_APIS, M.NET_WIN_HTTP_APIS, M.NET_WINSOCK_APIS)

    --- We fill in both x86 and x64 reserved registry key handles.
    self.pre_reg_key_handles['0x80000000'] = 'HKEY_CLASSES_ROOT'
    self.pre_reg_key_handles['0xffffffff80000000'] = 'HKEY_CLASSES_ROOT'

    self.pre_reg_key_handles['0x80000001'] = 'HKEY_CURRENT_USER'
    self.pre_reg_key_handles['0xffffffff80000001'] = 'HKEY_CURRENT_USER'

    self.pre_reg_key_handles['0x80000002'] = 'HKEY_LOCAL_MACHINE'
    self.pre_reg_key_handles['0xffffffff80000002'] = 'HKEY_LOCAL_MACHINE'

    self.pre_reg_key_handles['0x800000003'] = 'HKEY_USERS'
    self.pre_reg_key_handles['0xffffffff80000003'] = 'HKEY_USERS'

    self.pre_reg_key_handles['0x800000004'] = 'HKEY_PERFORMANCE_DATA'
    self.pre_reg_key_handles['0xffffffff80000004'] = 'HKEY_PERFORMANCE_DATA'

    self.pre_reg_key_handles['0x800000005'] = 'HKEY_CURRENT_CONFIG'
    self.pre_reg_key_handles['0xffffffff80000005'] = 'HKEY_CURRENT_CONFIG'

    self.pre_reg_key_handles['0x800000006'] = 'HKEY_DYN_DATA'
    self.pre_reg_key_handles['0xffffffff80000006'] = 'HKEY_DYN_DATA'

    self.pre_reg_key_handles['0x800000007'] = 'HKEY_CURRENT_USER_LOCAL_SETTINGS'
    self.pre_reg_key_handles['0xffffffff80000007'] = 'HKEY_CURRENT_USER_LOCAL_SETTINGS'

    self.pre_reg_key_handles['0x800000050'] = 'HKEY_PERFORMANCE_TEXT'
    self.pre_reg_key_handles['0xffffffff80000050'] = 'HKEY_PERFORMANCE_TEXT'

    self.pre_reg_key_handles['0x800000060'] = 'HKEY_PERFORMANCE_NLSTEXT'
    self.pre_reg_key_handles['0xffffffff80000060'] = 'HKEY_PERFORMANCE_NLSTEXT'

    self.Win32API = Win32API

end

M:init()
return M
