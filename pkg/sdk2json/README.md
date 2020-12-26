# Windows SDK API definition in JSON

sdk2json is a go package that parses the Windows SDK (Prototypes, structures, unions) to JSON format. 

Here is an example:

```json
{
 "advapi32.dll": {
  "ControlService": {
   "callconv": "WINAPI",
   "name": "ControlService",
   "retVal": "BOOL",
   "params": [
    {
     "anno": "_In_",
     "type": "SC_HANDLE",
     "name": "hService"
    },
    {
     "anno": "_In_",
     "type": "DWORD",
     "name": "dwControl"
    },
    {
     "anno": "_Out_",
     "type": "LPSERVICE_STATUS",
     "name": "lpServiceStatus"
    }
   ]
  },
```

The malware sandbox hooking module make use of this to implement a generic hook handler that does not require to implement a handler for each API we need to hook.


## Challenges

- Not a consistent way of defines functions prototypes or structs.