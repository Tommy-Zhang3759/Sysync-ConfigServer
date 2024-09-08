# Sysync ConfigServer

## Basic functions 
- Batch modify the Windows registry
- Batch modify Windows network configuration 
- WebUI
- Web client CLI
- Find and connect to the server according to the local settings.
- Listening to the broadcast and automatically get static IP and host names when the server address that is saved locally is unavailable. 
- Separate files/folders (NFS/SMB)
## Advanced features 
- Remote control (VNC)
- Screen monitoring 
## Development 
- Language
  - Server: Go (net APIs & Time wheel) & C++ (Authentication system)
  - Client: Python (Temporary use. Later use go to reconstruct the api part)
- Network communication: 
  - Format: JSON (or YAML, if there is necessity)
  - Encoding: UTF-8
  - Protocol: UDP
