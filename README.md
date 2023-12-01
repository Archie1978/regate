# Golang Remote Manager

Regate is a pure Golang implementation of the Microsoft RDP client and Ssh client.

I use adaptation code of tomatome/grdp ( and Sylvain Peyrefitte for JS ) for RDP and jsterminal dor SSH.


## Status

**The project is under development and not finished yet.**
Regate:
* [ ] Interface installation
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
* [ ] Standard SSH by HSM ( certificate x509 )

## Technologies

* Golang v19.0
* JsTerminal
* WebSocket

## Build ( step 1)
1. Build vue code
2. cd www/regate
3. yarn install
4. yarn run build
5. cd ../..

## Build Standalone ( step 2)
1. cd cmd/regate-standalone-user
2. go build -a

## Build Daemon mode bastion/mutiuser ( step 2)
1. cd cmd/regate-daemon/
2. go build -a 


## Release
Version: 0.0.0
Version init


## Plan dev Next
* Interface installation ( standalone )
* Use Regate standalone by user linux/Windows
* One binary HTML is into binary
* Use Regate multiuser (bastion) connexion LDAP

