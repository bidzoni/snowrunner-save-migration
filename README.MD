# Xbox PC Game Pass Snowrunner saves migration tool

This is a small app to convert Snowrunner's save files from Xbox Game Pass to Steam. 

Xbox stores save files not in `.cfg` format, so we need to rename or copy the files before copy than into a directory with saves in Steam.

## How to use
Locate your xbox save files
1. Open directory `%USERPROFILE%\AppData\Local\Packages\FocusHomeInteractiveSA.SnowRunnerWindows10_4hny5m903y3g0\`
2. Find subdir with name like `C4D962C08FBB42F08D364BEBB933C207`

In this directory locate container file with name like `container.52`. This file contains mapping rules for file names. 

Create a folder for migration, for example `C:\snow_migrate`

Run migration tool with command `snowrunner-save-migration.exe -container container.52 -source C:\Users\user\AppData\Local\Packages\FocusHomeInteractiveSA.SnowRunnerWindows10_4hny5m903y3g0\C4D962C08FBB42F08D364BEBB933C207 -dest C:\snow_migrate`

Now you have `.cfg` saves in `C:\snow_migrate` dir. Just copy them to the steam saves directory (something like `D:\Steam\userdata\42424242\1465360\remote`) and start the game. 
