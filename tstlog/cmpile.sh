

#!/bin/bash

 go build


sshpass -p bxyvrv1601 scp ./tstlog root@192.168.166.41:/opt/zap_http/
sshpass -p bxyvrv1601 scp ./blog4go_config.xml root@192.168.166.41:/opt/zap_http/
sshpass -p bxyvrv1601 scp ./seelog.xml root@192.168.166.41:/opt/zap_http/



