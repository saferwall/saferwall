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
    FILE_CREATE_APIs = {'CreateFileA', 'CreateFileW', 'CreateDirectory', 'CreateDirectoryExA', 'CreateDirectoryExW'},
    FILE_OPEN_APIs = {'OpenFile'},
    FILE_READ_APIs = {'ReadFile', 'ReadFileEx'},
    FILE_WRITE_APIs = {'WriteFile', 'WriteFileEx'},
    File_COPY_APIS = {'CopyFileA', 'CopyFileW', 'CopyFileExA', 'CopyFileExW'},
    File_MOVE_APIs = {'MoveFileA', 'MoveFileW', 'MoveFileWithProgressA', 'MoveFileWithProgressW'},
    FILE_DELETE_APIS = {'DeleteFileA', 'DeleteFileW'},

    --- Network APIs.
    NET_WIN_HTTP_APIS = {'WinHttpConnect'},
    NET_WIN_INET_APIS = {'InternetConnectA', 'InternetConnectW'},
    NET_WINSOCK_APIS = {'getaddrinfo', 'GetAddrInfoW', 'GetAddrInfoExA', 'GetAddrInfoExW'},

    -- Default port numbers for WinINet.
    DEFAULT_PORT = {INVALID = 0, FTP = 21, GOPHER = 70, HTTP = 80, HTTPS = 443, SOCKS = 1080},

    REG_SET_APIS = {},
    REGISTRY_APIS = {},
    FILE_APIS = {},
    NETWORK_APIS = {},

    --- Predefined registry key handles.
    pre_reg_key_handles = {},

    --- Handle types constants.
    HANDLE_TYPE = {REGISTRY = 1, FILE = 2, NETWORK = 3},

    -- Context table to this run.
    ctx = {

        -- Track Key Handle -> Object Name.
        reg_key_handles = {},
        file_handles = {},
        network_handles = {},

        -- Summarized events.
        events = {},

        -- Seen events.
        seen_events = {},

        -- Track created files.
        created_files = {}
    }
}

utils.table_merge(M.REG_SET_APIS, M.REG_SET_VALUE_APIS, M.REG_SET_KEY_VALUE_APIS)
utils.table_merge(M.REGISTRY_APIS, M.REG_CREATE_APIS, M.REG_OPEN_APIS, M.REG_SET_APIS, M.REG_DELETE_APIS)
utils.table_merge(M.FILE_APIS, M.FILE_CREATE_APIs, M.FILE_OPEN_APIs, M.FILE_READ_APIs, M.FILE_WRITE_APIs,
                  M.File_COPY_APIS, M.File_MOVE_APIs, M.FILE_DELETE_APIS)
utils.table_merge(M.NETWORK_APIS, M.NET_WIN_INET_APIS, M.NET_WIN_HTTP_APIS, M.NET_WINSOCK_APIS)

function M:init()

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

end

--- Evaluate an API trace against a set of rules.
---@param w32apis table<integer, Win32API> Table containing list of API objects.
function M:eval(w32apis)

    self.ctx.reg_key_handles = utils.table_copy(self.pre_reg_key_handles)

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

        local event = {proc_id = w32api.process_id}

        -- Summarize all registry operations.
        if utils.table_contains(self.REGISTRY_APIS, w32api.name) then
            -- hKey is either a handle returned by on of registry creation APIs;
            -- or it can be one of the predefined keys.
            local hKeyStr = w32api:param_value_from_name('hKey')

            -- lpSubKey is subkey of the key identified by the hKey parameter.
            local lpSubKey = w32api:param_value_from_name('lpSubKey')

            -- lpValueName is the name of the registry value whose data is to be updated.
            local lpValueName = w32api:param_value_from_name('lpValueName')

            -- phkResult is a pointer to a variable that receives a handle to the
            -- opened or created key.
            local phkResult = w32api:param_value_from_name('phkResult')

            event.type = 'registry'

            if utils.table_contains(self.REG_CREATE_APIS, w32api.name) then
                event.operation = 'create'
                event.path = self.ctx.reg_key_handles[hKeyStr] .. '\\' .. lpSubKey

                -- Save the mapping between the handle and its equivalent path.
                self.ctx.reg_key_handles[phkResult] = event.path
            elseif utils.table_contains(self.REG_OPEN_APIS, w32api.name) then
                event.operation = 'open'
                event.path = self.ctx.reg_key_handles[hKeyStr] .. '\\' .. lpSubKey

                -- Save the mapping between the handle and its equivalent path.
                self.ctx.reg_key_handles[phkResult] = event.path
            elseif utils.table_contains(self.REG_SET_APIS, w32api.name) then
                event.operation = 'write'

                if utils.table_contains(self.REG_SET_KEY_VALUE_APIS, w32api.name) then
                    event.path = self.ctx.reg_key_handles[hKeyStr] .. '\\' .. lpSubKey .. '\\\\' .. lpValueName
                else
                    event.path = self.ctx.reg_key_handles[hKeyStr] .. '\\\\' .. lpValueName
                end
            elseif utils.table_contains(self.REG_DELETE_APIS, w32api.name) then
                event.operation = 'delete'
                event.path = self.ctx.reg_key_handles[hKeyStr]

                if lpSubKey ~= nil then
                    event.path = event.path .. '\\' .. lpSubKey
                end

                if lpValueName ~= nil then
                    event.path = event.path .. '\\\\' .. lpValueName
                end
            end

            -- Summarize all registry operations.
        elseif utils.table_contains(self.FILE_APIS, w32api.name) then
            -- lpFileName points to the name of the file or device to be created or opened.
            local lpFileName = w32api:param_value_from_name('lpFileName')

            -- hFile represents a handle to the file or I/O device.
            local hFile = w32api:param_value_from_name('hFile')

            -- lpPathName points to the path of the directory to be created.
            local lpPathName = w32api:param_value_from_name('lpPathName')

            -- lpNewDirectory points to the path of the directory to be created.
            local lpNewDirectory = w32api:param_value_from_name('lpNewDirectory')

            -- lpExistingFileName points to the name of an existing file.
            local lpExistingFileName = w32api:param_value_from_name('lpExistingFileName')

            -- lpNewFileName points to the name of the new file.
            local lpNewFileName = w32api:param_value_from_name('lpNewFileName')

            -- Th return value of the API, which is a handle in the case of file APIs.
            local returnedHandle = w32api.return_value

            event.type = 'file'

            if utils.table_contains(self.FILE_CREATE_APIs, w32api.name) then
                event.operation = 'create'

                -- Either a file or a directory creation.
                if lpFileName ~= '' then
                    event.path = lpFileName
                else
                    -- The Ex version of create directory have a different param name.
                    if lpPathName ~= '' then
                        event.path = lpPathName
                    else
                        event.path = lpNewDirectory
                    end
                end

                -- Save the mapping between the handle and its equivalent path.
                self.ctx.file_handles[returnedHandle] = event.path

            elseif utils.table_contains(self.FILE_OPEN_APIs, w32api.name) then
                event.operation = 'open'
                event.path = lpFileName

                -- Save the mapping between the handle and its equivalent path.
                self.ctx.file_handles[returnedHandle] = event.path
            elseif utils.table_contains(self.FILE_READ_APIs, w32api.name) then
                event.operation = 'read'
                event.path = self:handle_to_object_name(hFile, self.HANDLE_TYPE.FILE)
            elseif utils.table_contains(self.FILE_WRITE_APIs, w32api.name) then
                event.operation = 'write'
                event.path = self:handle_to_object_name(hFile, self.HANDLE_TYPE.FILE)
            elseif utils.table_contains(self.FILE_DELETE_APIS, w32api.name) then
                event.operation = 'delete'
                event.path = lpFileName
            elseif utils.table_contains(self.File_COPY_APIS, w32api.name) then
                event.operation = 'copy'
                event.path = lpExistingFileName .. '->' .. lpNewFileName
            elseif utils.table_contains(self.File_MOVE_APIs, w32api.name) then
                event.operation = 'move'
                event.path = lpExistingFileName .. '->' .. lpNewFileName
            end
        elseif utils.table_contains(self.NETWORK_APIS, w32api.name) then
            event.type = 'network'
            if utils.table_contains(self.NET_WINSOCK_APIS, w32api.name) then
                -- pNodeName contains a host (node) name or a numeric host address string.
                local pNodeName = w32api:param_value_from_name('pNodeName')

                -- pName contains a host (node) name or a numeric host address string.
                local pName = w32api:param_value_from_name('pName')

                -- pServiceName contains either a service name or port number represented as a string.
                local pServiceName = w32api:param_value_from_name('pServiceName')
                if pNodeName ~= nil then
                    event.path = pNodeName
                else
                    event.path = pName
                end

                event.path = event.path .. ':' .. pServiceName
                event.operation = 'socket'
            else
                -- lpszServerName specifies the host name of an Internet server.
                local lpszServerName = w32api:param_value_from_name('lpszServerName')

                -- nServerPort represents the TCP/IP port on the server.
                local nServerPort = w32api:param_value_from_name('nServerPort')

                -- pswzServerName contains the host name of an HTTP server.
                local pswzServerName = w32api:param_value_from_name('pswzServerName')

                -- The return value of the API, which is a handle in the case of network APIs.
                local returnedHandle = w32api.return_value

                if utils.table_contains(self.NET_WIN_HTTP_APIS, w32api.name) then
                    event.path = pswzServerName
                elseif utils.table_contains(self.NET_WIN_INET_APIS, w32api.name) then
                    event.path = lpszServerName
                end

                local server_port = tonumber(nServerPort)
                if server_port == self.DEFAULT_PORT.INVALID then
                    -- Uses the default port for the service specified by dwService.
                    local dwService = w32api:param_value_from_name('dwService')
                    local svc_port = tonumber(dwService)
                    if svc_port == 0x1 then
                        event.operation = 'ftp'
                        server_port = 21
                    elseif svc_port == 0x2 then
                        event.operation = 'gopher'
                        server_port = 70
                    elseif svc_port == 0x3 then
                        event.operation = 'http'
                        server_port = 80
                    end
                elseif server_port == self.DEFAULT_PORT.FTP then
                    event.operation = 'ftp'
                elseif server_port == self.DEFAULT_PORT.HTTP then
                    event.operation = 'http'
                elseif server_port == self.DEFAULT_PORT.HTTPS then
                    event.operation = 'https'
                elseif server_port == self.DEFAULT_PORT.SOCKS then
                    event.operation = 'socks'
                end

                event.path = event.path .. ':' .. server_port

                -- Save the mapping between the handle and its equivalent path.
                self.ctx.network_handles[returnedHandle] = event.path
            end
        end

        if event.path ~= nil then
            local key = event.proc_id .. event.path .. event.operation
            if self.ctx.seen_events[key] == nil then
                table.insert(self.ctx.events, event)
                self.ctx.seen_events[key] = true
            end

        end

    end

    return nil

    -- create a table that holdes all creates files.
end

function M:handle_to_object_name(handle, type)
    if type == self.HANDLE_TYPE.FILE then
        return self.ctx.file_handles[handle]
    elseif type == self.HANDLE_TYPE.REGISTRY then
        return self.ctx.reg_key_handles[handle]
    elseif type == self.HANDLE_TYPE.NETWORK then
        return self.ctx.network_handles
    end
end

M:init()
return M
