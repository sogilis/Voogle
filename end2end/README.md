# End to End tests

With Voogle running on you machine, execute `make` and the tests must pass.

## /!\ Warning 
For now, the minio MUST be EMPTY before running them.

So we need to determine the proper way to handle this.

Maybe counting the number of elements returned to the list is not the right solution. And we should only check that the request succeed

