# Disk Image Backup over WEB

You can take disk image over with Dibow over web. You don't need any space at
 your server.



### Recommended Use Cases
- Mirror your disks at remote location.
- Replicate your server.
- Transfer server between different companies.

**NOTE:** If you want to take a root system image, do not use this software at
live linux partition. For the taking root partition image, use Dibow at rescue disk or mode. Otherwise your image will be corrupted.

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
# Run service
./dibow
```


## Program Options
For a now just allow only one option.

### listen-addr

Program normally starts with port 80 but if you want to use a different you can
use a --listen-addr argument to set different port or address

```bash
./dibow --listen-addr :8443
```
