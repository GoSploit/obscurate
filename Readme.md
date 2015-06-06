Obscurate
=========

Obscurate is a simple library designed to reduce fingerprintable code in GoSploit payloads.

Obscuration is similar to encryption, but is not designed to be secure.  The primary goal of obscurate is 2-fold:

1. Make data harder to identify without non-trivial levels of work.
2. Reduce fingerprintable code.

Fingerprintable code is code that doesn't change between "variants" of a payload.

In standard encryption, you usually end up with 3 things needed to get back the original data:

1. The encrypted data
2. The decryption key
3. The decryption algorithm

The issue is, that algorithm is a constant, that could possibly be fingerprinted to identify payloads.

Obscuration combines the decryption "key" and "algorithm".
By doing so, each obscuration of the same chunk of data will have a different "encrypted data" and algorithm.

Note: obscuration, and the code produced by the obscurate library, is not cryptographically secure.
If you want to prevent a 3rd party from reading your data, use a cryptographically secure library.
