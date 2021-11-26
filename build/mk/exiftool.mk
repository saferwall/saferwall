EXIF_VER = 12.36
exiftool-install: # Install ExifTool
	sudo apt-get -qq update
	wget https://exiftool.org/Image-ExifTool-$(EXIF_VER).tar.gz
	gzip -dc Image-ExifTool-$(EXIF_VER).tar.gz | tar -xf -
	cd Image-ExifTool-$(EXIF_VER) \
		&& perl Makefile.PL \
		&& make test \
		&& sudo make install
	rm Image-ExifTool-$(EXIF_VER).tar.gz
	rm -r Image-ExifTool-$(EXIF_VER)