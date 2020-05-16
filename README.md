# Disk Image Backup over WEB

You can take disk image over with Dibow over web. You don't need any space at your server for the taking backup.

Dibow project is created for taking disk image backup to remote locations or transferring server images into another remote server.

### Recommended Use Cases
- Mirror your disks at remote location.
- Replicate your server.
- Transfer server between different companies.

**NOTE:** If you want to take a root system image, do not use this software at live linux partition. For the taking root partition image, use Dibow at rescue disk or mode. Otherwise your image will be corrupted.

## Installing dibow

### Compiling from source
For the Compiling you need golang at your machine.  
You can get source from github or gitlab. If you are using IPv6 only network you can use gitlab instead of github.
```bash
# Getting source
git@gitlab.com:ahmetozer/dibow.git
# Enter source directory
cd dibow
# Build from source
go build
```

### Get from gitlab

You can get binary from gitlab.
```bash
wget "https://gitlab.com/ahmetozer/dibow/-/jobs/artifacts/master/download?job=compile" -O dibow.zip
unzip dibow.zip
chmod +x dibow
```

## Program Modes

### Server
This option is serve your disk into web service.
You can list and download your disk images on web.

Default username is `root` .  
A random password will be displayed on the terminal.

#### Arguments
- `--listen-addr` Program normally starts with port 443 but if you want to use a different you can use a --listen-addr argument to set different port or address.
```bash
./dibow server --listen-addr :8443
```

<!---  
### Client
You can get or write disk image to on remote systems.
#### Arguments


- `--save` Save image to your pc.
```bash
./dibow client --save --url https://example.com/image/dev/sda
# To save different location and different name
./dibow client --save /root/backup/old_server.img --url https://example.com/image/dev/sda
```
- `--write ` Write remote image to disk.
```bash
./dibow client --write /dev/sda --url https://example.com/image/dev/sda
```
-->
