go get -u github.com/c-darwin/dcoin-go-tmp
cd ../
GOARCH=amd64  CGO_ENABLED=1  go build -o make_dmg/dcoin.app/Contents/MacOs/dcoinbin
cd make_dmg
zip -r dcoin_osx64.zip dcoin.app
./make_dmg.sh -b background.png -i logo-big.icns -s "480:540" -c 240:400:240:200 -n dcoin64 "dcoin.app"
cd ../
GOARCH=386  CGO_ENABLED=1  go build -o make_dmg/dcoin.app/Contents/MacOs/dcoinbin
cd make_dmg
zip -r dcoin_osx32.zip dcoin.app
./make_dmg.sh -b background.png -i logo-big.icns -s "480:540" -c 240:400:240:200 -n dcoin32 "dcoin.app"
 

