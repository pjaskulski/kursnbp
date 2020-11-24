version='0.3.1'

pack_bsd=kursnbp_v${version}_freebsd_amd64.tar.gz
pack_lin=kursnbp_v${version}_linux_amd64.tar.gz
pack_win=kursnbp_v${version}_windows_amd64.zip
pack_mac=kursnbp_v${version}_macos_amd64.tar.gz

# FreeBSD
echo "FreeBSD..."
FILE=releases/freebsd/$pack_bsd
if [ -f "$FILE" ]; then
    rm $FILE
fi
cd ./builds/freebsd
tar czvf $pack_bsd kursnbp
mv $pack_bsd ../../releases/freebsd/$pack_bsd
cd ../..

# Linux
FILE=releases/linux/$pack_lin
if [ -f "$FILE" ]; then
    rm $FILE
fi
cd ./builds/linux
tar czvf $pack_lin kursnbp
mv $pack_lin ../../releases/linux/$pack_lin
cd ../..

# Windows
FILE=releases/windows/$pack_win
if [ -f "$FILE" ]; then
    rm $FILE
fi
cd ./builds/windows
zip $pack_win kursnbp.exe
mv $pack_win ../../releases/windows/$pack_win
cd ../..

# MacOS
FILE=releases/macos/$pack_mac
if [ -f "$FILE" ]; then
    rm $FILE
fi
cd ./builds/macos
tar czvf $pack_mac kursnbp
mv $pack_mac ../../releases/macos/$pack_mac
cd ../..

echo "Done!"