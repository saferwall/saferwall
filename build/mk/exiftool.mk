EXIF_VER = 12.62
exiftool-install: # Install ExifTool
	wget https://exiftool.org/Image-ExifTool-$(EXIF_VER).tar.gz
	gzip -dc Image-ExifTool-$(EXIF_VER).tar.gz | tar -xf -
	cd Image-ExifTool-$(EXIF_VER) \
		&& perl Makefile.PL \
		&& make test \
		&& sudo make install
	cd Image-ExifTool-$(EXIF_VER) && sudo cp -r exiftool lib /usr/local/bin
	rm Image-ExifTool-$(EXIF_VER).tar.gz
	rm -rf Image-ExifTool-$(EXIF_VER)
