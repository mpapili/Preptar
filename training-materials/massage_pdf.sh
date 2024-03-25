#! /bin/bash

# this is combining a few "sed" commands to try and turn most PDF's into a
# semi-usable text file

echo "what is the name of the pdf (without .pdf)?"
read filename

pdftotext -layout ${filename}.pdf ${filename}.txt
filename=${filename}.txt


# turn line breaks that aren't after punctuation into spaces
sed -i ':a;N;$!ba;s/\([^.!?]\)\n/\1 /g' $filename
# turn all instances of multiple spaces into single spaces
sed -i 's/ \{2,\}/ /g' $filename

# remove non-ascii characters:
sed -i 's/ \{2,\}/ /g' $filename

echo "Done, new text file to use is $filename