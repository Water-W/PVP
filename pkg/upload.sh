#! /bin/sh
## author
## this is a shell script

ips=("39.104.200.8" "39.104.200.4" "39.104.204.99" "39.104.200.94" "39.104.202.92" "39.104.203.170" "39.104.205.191" "39.104.204.231" "39.104.203.26" "39.104.207.231")
password="dev@BDW"


for ip in ${ips[@]};
do
    set time 20
    echo 'scp ../pvp dev@'$ip':/home/dev/hu'
    # expect {
    # "*yes/no*"    {send   "yes\r";exp_continue}
    # "*password" {send   "$password\r"}
    # }
    # expect eof
done

scp /home/wangyumin/go/src/PVP/PVP/pvp dev@39.104.200.8:/home/dev/hu
scp /home/wangyumin/go/src/PVP/PVP/pvp dev@39.104.200.4:/home/dev/hu
scp /home/wangyumin/go/src/PVP/PVP/pvp dev@39.104.204.99:/home/dev/hu
scp /home/wangyumin/go/src/PVP/PVP/pvp dev@39.104.200.94:/home/dev/hu
scp /home/wangyumin/go/src/PVP/PVP/pvp dev@39.104.202.92:/home/dev/hu
scp /home/wangyumin/go/src/PVP/PVP/pvp dev@39.104.203.170:/home/dev/hu
scp /home/wangyumin/go/src/PVP/PVP/pvp dev@39.104.205.191:/home/dev/hu
scp /home/wangyumin/go/src/PVP/PVP/pvp dev@39.104.204.231:/home/dev/hu
scp /home/wangyumin/go/src/PVP/PVP/pvp dev@39.104.203.26:/home/dev/hu
scp /home/wangyumin/go/src/PVP/PVP/pvp dev@39.104.207.231:/home/dev/hu


scp -r /home/wangyumin/go/src/PVP/PVP/frontend dev@39.104.200.8:/home/dev/hu
scp -r /home/wangyumin/go/src/PVP/PVP/frontend dev@39.104.200.4:/home/dev/hu
scp -r /home/wangyumin/go/src/PVP/PVP/frontend dev@39.104.204.99:/home/dev/hu
scp -r /home/wangyumin/go/src/PVP/PVP/frontend dev@39.104.200.94:/home/dev/hu
scp -r /home/wangyumin/go/src/PVP/PVP/frontend dev@39.104.202.92:/home/dev/hu
scp -r /home/wangyumin/go/src/PVP/PVP/frontend dev@39.104.203.170:/home/dev/hu
scp -r /home/wangyumin/go/src/PVP/PVP/frontend dev@39.104.205.191:/home/dev/hu
scp -r /home/wangyumin/go/src/PVP/PVP/frontend dev@39.104.204.231:/home/dev/hu
scp -r /home/wangyumin/go/src/PVP/PVP/frontend dev@39.104.203.26:/home/dev/hu
scp -r /home/wangyumin/go/src/PVP/PVP/frontend dev@39.104.207.231:/home/dev/hu