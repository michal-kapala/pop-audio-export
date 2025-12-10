# Prince of Persia 2008 audio exporter
A rudimentary `.forge` parser for exporting audio from the game files.

## Files
Possible exported formats are:
- `.bao` - [Binary Audio Object](https://github.com/vgmstream/vgmstream/blob/ed976476635829ecb23b26b074a0c03ecabd0f7a/src/meta/ubi_bao.c) file for Ubisoft's proprietary Dare sound engine
- `.data` - a non-BAO asset/metadata file
- `.ogg` - raw audio extracted from BAO

The game uses OGG/Vorbis in most BAOs (~90%), for other formats/codecs use [vgmstream](https://github.com/vgmstream/vgmstream).

## Credits
Special thanks to [@forge-master](https://github.com/forge-master), [@Kamzik123](https://github.com/Kamzik123) and [Blacksmith](https://github.com/theawesomecoder61/Blacksmith) contributors for help with forge/bao formats.
