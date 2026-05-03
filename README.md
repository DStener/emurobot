# emurobot

Emulation is performed exclusively for the devices defined in the configuration file. Only devices listed under the devices key will be emulated. Refer to the example configuration below for the required format

```yaml
version: 1.0.0
devices: 
  - input: /dev/ttyUSB0
    output: /dev/ttyUSB0_VIRT
    speed: 9600
    size: 8
  # ...
```

## Running 

Build control tool

```sh
go build -o build/emurobot ./cli
```

### For simulate real environment

1. Up docker-compose
```sh
docker compose -f docker-compose.emu_record.yml up
```

2. Send start record command 
```sh
EMU_SIGNAL_ADDRESS=http://localhost:7000 ./build/emurobot -rec start
```

3. Wait...

4. Send stop record command
```sh
EMU_SIGNAL_ADDRESS=http://localhost:7000 ./build/emurobot -rec stop
```

Now in dumps/{DATA}/{TIME}_dumps_SPs.json has written logs  

### For emulate environment (play records)

1. Up docker-compose
```sh
docker compose -f docker-compose.emu_play.yml up
```

Now the devices specified in the configuration file will exist in the system 

2. Send play record command with path to log file 
```sh
EMU_SIGNAL_ADDRESS=http://localhost:7000 ./build/emurobot -play /app/dumps/03.05.2026/15_06_39_dumps_SPs.json
```

3. Check devices 