# The first 1024 bytes of this file need to contain only ASCII characters.
# After the first 1024 bytes, then there should be a non-ASCII character.
#
# Explanation:
# We use golang.org/x/net/html/charset.DetectEncoding to guess file encoding.
# DetectEncoding checks the first 1024 bytes of a file.
# If it can't determine the encoding and saw no non-ASCII characters,
# it declares the file to have windows-1252 encoding.
# This mangles the submitted file if it should have been UTF-8.
# We test to make sure we use UTF-8 for such files, instead of windows-1252.

lipsum = """
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam condimentum vitae
ipsum eget tempor. Morbi sed ex quis orci vulputate cursus quis non massa.
Vestibulum quam nibh, elementum in justo in, venenatis tristique nisl. Morbi
sagittis elit id velit ultricies, sed rutrum augue posuere. Donec nec nulla nec
eros fringilla pellentesque. Duis at dictum justo. Nunc ut magna felis. Aliquam
volutpat, lectus et molestie porttitor, est orci malesuada erat, ac pretium
eros ligula vel erat. Nullam venenatis dui eget sapien semper lobortis. Aenean
ac eros eget neque porta auctor in nec erat. Phasellus ac nulla ac turpis
porttitor auctor. Etiam eget posuere diam, ac feugiat lacus. Curabitur ornare
justo ut nulla congue, vitae posuere erat venenatis. Aliquam pulvinar eleifend
faucibus.

Etiam justo sem, faucibus malesuada purus a, ultrices efficitur ex.
Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac
turpis egestas. Duis maximus dapibus mattis. Quisque sem ex, convallis eu
ultricies posuere.
"""

# 👍
