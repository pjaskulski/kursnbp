version='0.3.7'

pack_bsd=kursnbp_v${version}_freebsd_amd64.tar.gz
pack_lin=kursnbp_v${version}_linux_amd64.tar.gz
pack_win=kursnbp_v${version}_windows_amd64.zip
pack_mac=kursnbp_v${version}_macos_amd64.tar.gz

# FreeBSD
echo "FreeBSD..."
if [ ! -n "$(find releases/freebsd/ -prune -empty 2>/dev/null)" ]
then
  rm releases/freebsd/kursnbp_*
fi
cd ./builds/freebsd
tar czvf $pack_bsd kursnbp
mv $pack_bsd ../../releases/freebsd/$pack_bsd
cd ../..

# Linux
if [ ! -n "$(find releases/linux/ -prune -empty 2>/dev/null)" ]
then
  rm releases/linux/kursnbp_*
fi
cd ./builds/linux
tar czvf $pack_lin kursnbp
mv $pack_lin ../../releases/linux/$pack_lin
cd ../..

# Windows
if [ ! -n "$(find releases/windows/ -prune -empty 2>/dev/null)" ]
then
  rm releases/windows/kursnbp_*
fi
cd ./builds/windows
zip $pack_win kursnbp.exe
mv $pack_win ../../releases/windows/$pack_win
cd ../..

# MacOS
if [ ! -n "$(find releases/macos/ -prune -empty 2>/dev/null)" ]
then
  rm releases/macos/kursnbp_*
fi
cd ./builds/macos
tar czvf $pack_mac kursnbp
mv $pack_mac ../../releases/macos/$pack_mac
cd ../..

echo "Done!"