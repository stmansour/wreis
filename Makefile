DIRS=util db csv ws server

wreis:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done

test:
	for dir in $(DIRS); do make -C $$dir test;done

package:

all: clean wreis test
