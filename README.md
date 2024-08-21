# Device Service 

This is a service to store and access data related to devices.

### Build and run (locally)
To build the binary, run:
```ssh
make build
```
the excutable will be stored in ./bin/svc

Or simply run 
```ssh
make
```
To build and run the binary 

### API 
> All paths start with: /api/vi  
> e.g. GET http://localhost:8000/api/v1/devices, to list all devices

| Method | Path                     | Description                                       |
|--------|--------------------------|---------------------------------------------------|
| GET    | /devices                 | List all devices                                  |
| GET    | /devices?s=brand:< value > | List all devices with the specified brand (value) |
| GET    | /devices/< id >            | Get a device with the specified id                |
| POST   | /devices                 | Post a device giving a name and a brand           |
| PUT    | /devices/< id >            | Update partial or full a device given the id      |
| DELETE | /devices/< id >            | Delete a device with the specified id             |

### Test
To test the source code, run:   
```ssh
make test
```
The important package that needed to be test are:
- controllers
- repositories

Alos, the coverage is printed along side each package

### Containerize
First to build the image, run:
```ssh
make build-image
```

And then:
```ssh
make run-image
```
This will start a container running the device service

To stop and delete the image, run:
```ssh
make rmi
```