Just print it!
--------------

I hate printer drivers.  And printer servers.  And printers.

Since it's impossible to have all three of these things working on all
machines at the same time, I'm punting on the problem.  New plan:

  * Set up the printer on a single machine and be careful never to
    touch it again

  * Run this little server to let people on other machines upload
    their files via a web browser and have them printed

Runs on Linux or OS X as follows:

     ./just_print_it.linux 8080 'Brother_HL_2130_series' 2>&1 >> logfile.log

     ./just_print_it.osx 8080 'Brother_HL_2130_series' 2>&1 >> logfile.log

If you have Go 1.8+ installed, builds with `./build.sh`.  Totally
non-standard build process, but I like that!
