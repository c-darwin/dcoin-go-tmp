#echo $1;
#ARCH="amd64"
#if [ $1 = "32" ];
#then
#  ARCH="386"
#fi

#echo $ARCH

go get -u github.com/c-darwin/dcoin-go-tmp
GOARCH=amd64  CGO_ENABLED=1  go build -o dcoin.app/Contents/MacOs/dcoinbin
./make_dmg.sh -b background.png -i logo-big.icns -s "480:540" -c 240:400:240:200 -n dcoin64 "dcoin.app"
GOARCH=386  CGO_ENABLED=1  go build -o dcoin.app/Contents/MacOs/dcoinbin
./make_dmg.sh -b background.png -i logo-big.icns -s "480:540" -c 240:400:240:200 -n dcoin32 "dcoin.app"
 

