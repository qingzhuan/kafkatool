p=$(ps -ef |grep kafkatool |grep -v grep|grep kafkatool)
if [[ -n "$p" ]];then
    echo $p |awk '{print $2}' |xargs kill
fi
go build
nohup ./kafkatool &