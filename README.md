## duperig

Find copied and modified (*duped* and *rigged*) files between two related folder hierarchies which are contained in git repositories. Prints current file SHA-256 along with its latest matching commit hash **as seen in the git repository of the first folder specified in the argument list** (if available).

### Requirements

Some version of git which isn't ancient. Tested with git `2.17.1`.

### Installation

```
go get -u github.com/MMulthaupt/duperig
```

### Usage

`duperig path/to/folder/from/which/files/were/copied path/to/folder/to/which/files/were/copied`

### Example

```
duperig projects/thingamajig_base/src/main/java/com/pany projects/thingamajig_special/src/main/java/com/pany
DIFF: services/Foo.java: 7418a5b686 (Commit: a7e5433603) vs dd232c048b (Commit: 119bb3cb41)
DIFF: data/Result.java: 711528c9c2 (Commit: a70d2cc30f) vs 2503ca3aeb (Commit: 47a05b7ada)
DIFF: data/Data.java: 123456789ab (Commit: 98765434567) vs 2637485985 (NO MATCHING COMMIT)
DUPE: mail/MailClient.java @ 3ec54907a4
DUPE: save/Database.java @ 5c346980d6
```

Folder structure at `projects/thingamajig_special/src/main/java/com/pany` has 5 files with paths coninciding with 5 other files under `projects/thingamajig_base/src/main/java/com/pany`. Out of those 5 files, 2 are exact duplicates. The remaining 3 files differ. From the differing files, 2 have commit hashes in `projects/thingamajig_base`. However, the file `projects/thingamajig_base/src/main/java/com/data/Data.java` has changes unique to `projects/thingamajig_special`.

### Status

Working, but slow and unoptimized. Not suitable for projects above 1 GB.

### TODO

See issues.
