## duperig

Find copied and modified (*duped* and *rigged*) files between two related folder hierarchies which are contained in git repositories. Prints file SHA-256 along with commit hash (if available).

### Requirements

Some version of git which isn't ancient. Tested with git `2.17.1`.

### Usage

```
go get -u github.com/MMulthaupt/duperig`
duperig projects/thingamajig_base/src/main/java/com/pany projects/thingamajig_special/src/main/java/com/pany
DIFF: services/Foo.java: 7418a5b686 (Commit: a7e5433603) vs dd232c048b (Commit: 119bb3cb41)
DIFF: data/Result.java: 711528c9c2 (Commit: a70d2cc30f) vs 2503ca3aeb (Commit: 47a05b7ada)
DIFF: data/Data.java: 123456789ab (NO MATCHING COMMIT) vs 2637485985 (Commit: 98765434567)
DUPE: mail/MailClient.java @ 3ec54907a4
DUPE: save/Database.java @ 5c346980d6
```

### Status

Working, but slow and unoptimized. Not suitable for projects above 1 GB.
