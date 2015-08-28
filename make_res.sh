echo "######## generate R.java ########"
aapt package -v -f -J /home/z/go-projects/src/github.com/c-darwin/dcoin-go-tmp/ -S /home/z/go-projects/src/github.com/c-darwin/dcoin-go-tmp/res/ -M /home/z/go-projects/src/github.com/c-darwin/dcoin-go-tmp/AndroidManifest.xml -I /home/z/android-sdk-linux/platforms/android-22/android.jar
mv R.java /home/z/go-projects/src/github.com/c-darwin/dcoin-go-tmp/R/org/golang/app/
echo "######## generate R.jar ########"
cd R
jar cfv /home/z/go-projects/src/github.com/c-darwin/dcoin-go-tmp/R.jar .
cd ../
echo "######## generate unsigned.apk ########"
aapt package -v -f -J /home/z/go-projects/src/github.com/c-darwin/dcoin-go-tmp/ -S /home/z/go-projects/src/github.com/c-darwin/dcoin-go-tmp/res/ -M /home/z/go-projects/src/github.com/c-darwin/dcoin-go-tmp/AndroidManifest.xml -I /home/z/android-sdk-linux/platforms/android-22/android.jar -F unsigned.apk
echo "######## extract resources.arsc ########"
unzip unsigned.apk -d apk
mv apk/resources.arsc .
rm -rf apk unsigned.apk