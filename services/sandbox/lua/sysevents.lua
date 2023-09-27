local utils = require 'utils'
local w32 = require 'w32api'

local M = {
    EVENT_TYPE = {REGISTRY = 'registry', FILE = 'file', NETWORK = 'network'},
    OP = {
        CREATE = 'create',
        OPEN = 'open',
        READ = 'read',
        WRITE = 'write',
        MOVE = 'move',
        COPY = 'copy',
        DELETE = 'delete'
    }
}

---@class SystemEvent
---@field proc_id string Process identifier responsible for generating the event.
---@field type string Type of the system event.
---@field path string Path of the system event.
---@field operation string  The operation requested over the above `path` field.
local SystemEvent = {}

function SystemEvent:new(o, proc_id, type, path, operation)
    o = o or {}
    setmetatable(o, self)
    self.__index = self
    o.proc_id = proc_id
    o.type = type
    o.path = path
    o.operation = operation
    return o
end

---@param ctx Context Contains specific information about this run.
---@param w32api Win32API Timestamp when the API was executed
function M:summarize(ctx, w32api)
    local event = SystemEvent:new(nil, w32api.process_id)

    if utils.table_contains(w32.REGISTRY_APIS, w32api.name) then
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

        event.type = self.EVENT_TYPE.REGISTRY

        if utils.table_contains(w32.REG_CREATE_APIS, w32api.name) then
            event.operation = self.OP.CREATE
            event.path = ctx.reg_key_handles[hKeyStr] .. '\\' .. lpSubKey

            -- Save the mapping between the handle and its equivalent path.
            ctx.reg_key_handles[phkResult] = event.path
        elseif utils.table_contains(w32.REG_OPEN_APIS, w32api.name) then
            event.operation = self.OP.OPEN
            event.path = ctx.reg_key_handles[hKeyStr] .. '\\' .. lpSubKey

            -- Save the mapping between the handle and its equivalent path.
            ctx.reg_key_handles[phkResult] = event.path
        elseif utils.table_contains(w32.REG_SET_APIS, w32api.name) then
            event.operation = self.OP.WRITE

            if utils.table_contains(w32.REG_SET_KEY_VALUE_APIS, w32api.name) then
                event.path = ctx.reg_key_handles[hKeyStr] .. '\\' .. lpSubKey .. '\\\\' .. lpValueName
            else
                event.path = ctx.reg_key_handles[hKeyStr] .. '\\\\' .. lpValueName
            end
        elseif utils.table_contains(w32.REG_DELETE_APIS, w32api.name) then
            event.operation = self.OP.DELETE
            event.path = ctx.reg_key_handles[hKeyStr]

            if lpSubKey ~= nil then
                event.path = event.path .. '\\' .. lpSubKey
            end

            if lpValueName ~= nil then
                event.path = event.path .. '\\\\' .. lpValueName
            end
        end

        -- Summarize all registry operations.
    elseif utils.table_contains(w32.FILE_APIS, w32api.name) then
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

        event.type = self.EVENT_TYPE.FILE

        if utils.table_contains(w32.FILE_CREATE_APIS, w32api.name) then
            event.operation = self.OP.CREATE

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
            ctx.file_handles[returnedHandle] = event.path

        elseif utils.table_contains(w32.FILE_OPEN_APIS, w32api.name) then
            event.operation = self.OP.OPEN
            event.path = lpFileName

            -- Save the mapping between the handle and its equivalent path.
            ctx.file_handles[returnedHandle] = event.path
        elseif utils.table_contains(w32.FILE_READ_APIS, w32api.name) then
            event.operation = self.OP.READ
            event.path = w32:handle_to_object_name(ctx, hFile, w32.HANDLE_TYPE.FILE)
        elseif utils.table_contains(w32.FILE_WRITE_APIS, w32api.name) then
            event.operation = self.OP.WRITE
            event.path = w32:handle_to_object_name(ctx, hFile, w32.HANDLE_TYPE.FILE)
        elseif utils.table_contains(w32.FILE_DELETE_APIS, w32api.name) then
            event.operation = self.OP.DELETE
            event.path = lpFileName
        elseif utils.table_contains(w32.File_COPY_APIS, w32api.name) then
            event.operation = self.OP.COPY
            event.path = lpExistingFileName .. '->' .. lpNewFileName
        elseif utils.table_contains(w32.File_MOVE_APIS, w32api.name) then
            event.operation = self.OP.MOVE
            event.path = lpExistingFileName .. '->' .. lpNewFileName
        end
    elseif utils.table_contains(w32.NETWORK_APIS, w32api.name) then

        event.type = self.EVENT_TYPE.NETWORK

        if utils.table_contains(w32.NET_WINSOCK_APIS, w32api.name) then
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

            if utils.table_contains(w32.NET_WIN_HTTP_APIS, w32api.name) then
                event.path = pswzServerName
            elseif utils.table_contains(w32.NET_WIN_INET_APIS, w32api.name) then
                event.path = lpszServerName
            end

            local server_port = tonumber(nServerPort)
            if server_port == w32.DEFAULT_PORT.INVALID then
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
            elseif server_port == w32.DEFAULT_PORT.FTP then
                event.operation = 'ftp'
            elseif server_port == w32.DEFAULT_PORT.HTTP then
                event.operation = 'http'
            elseif server_port == w32.DEFAULT_PORT.HTTPS then
                event.operation = 'https'
            elseif server_port == w32.DEFAULT_PORT.SOCKS then
                event.operation = 'socks'
            end

            event.path = event.path .. ':' .. server_port

            -- Save the mapping between the handle and its equivalent path.
            ctx.network_handles[returnedHandle] = event.path
        end
    end

    if event.path ~= nil then
        local key = event.proc_id .. event.path .. event.operation
        if ctx.seen_events[key] == nil then
            table.insert(ctx.events, event)
            ctx.seen_events[key] = true
            return event
        end
    end

    return nil

end

return M
