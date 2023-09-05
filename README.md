# Golang Remote Manager

Regate is a pure Golang implementation of the Microsoft RDP client and Ssh client.

I use adaptation code of tomatome/grdp ( and Sylvain Peyrefitte for JS ) for RDP and jsterminal dor SSH.


## Status

**The project is under development and not finished yet.**
Regate:
* [ ] Interface (unfinished)
* [ ] Ssl configuration ( not yet)
* [ ] Authentification interface ( not yet)
* [ ] Security of password ( not yet)
* [ ] Administration account
* [ ] Single binary


RDP:
* [x] Standard RDP Authentication
* [x] SSL Authentication
* [x] NTLMv2 Authentication
* [ ] Windows Clipboard
* [ ] RDP Client(ugly)
* [ ] VNC Client(unfinished)

SSH:
* [x] Standard SSH by password
* [ ] Standard SSH by key
* [ ] Standard SSH by HSM

## Technologies

* Golang v19.0
* JsTerminal
* WebSocket

## Build
1. Build vue code
2. cd www/regate
3. npm run build
4. cd ../..
5. go build

