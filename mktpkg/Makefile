DIRS=js

mktpkg:
	for dir in $(DIRS); do cd $$dir;make;done
	echo "Completed"

test:
	echo "Completed"

clean:
	for dir in $(DIRS); do cd $$dir;make clean;done

package:
	for dir in $(DIRS); do cd $$dir;make package;done
	echo "Completed"
