sherlock
========

__this code is a proof-of-concept__

Sherlock helps rank "important" words in a given corpus. It tries to do this in as many unsupervised methods as possible; meaning, it tries to make as many assumptions about a word as it can and provides you the relevant scores to continue the classification process.


## Currently Implemented
* TF-IDF
* A variant of k-NN; using the distance of the next character in a word in the known alphabet (closer = improvement, further = penalty)
* Generic filtering:
  * Word cannot be empty
  * Word must be at least 4 characters long
  * Words are penalized for having characters not in the known alphabet
