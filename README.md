# go-acc

Compute CRC32 and ED2K hashes for files.

Default behavior is to simply compute and verify CRC32 hashes found in
filenames.  If no hash is found, one can be added with the `--add` option.

The `--add-delim` option sets the character to use to separate the original
filename from the added CRC32 hash.  This defaults to an empty string (ie, no
delimiter).  The added hash is always surrounded by square brackets.

Generating full ED2K links is done with the `--ed2k option`.  If this option is
passed, the files will not be CRC32 hashed.  ED2K hashes cannot
be added to filenames.

    Usage: main [--add] [--add-delim ADD-DELIM] [--ed2k] INPUTFILES [INPUTFILES ...]

    Positional arguments:
      INPUTFILES             Input files

    Options:
      --add, -a              Add the calculated hash to the filename if none is found
      --add-delim ADD-DELIM, -d ADD-DELIM
                             Character to use before the added hash
      --ed2k, -e             Print ED2K links
      --help, -h             display this help and exit
