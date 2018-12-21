

#!/bin/bash

GOOS=linux GOARCH=mipsle go build

sshpass -p loongson scp ./tstlog    loongson@192.168.166.247:/home/loongson/tstpath/
#sshpass -p bxyvrv1601 scp ./blog4go_config.xml root@192.168.166.247:/opt/loongson/tstpath/



