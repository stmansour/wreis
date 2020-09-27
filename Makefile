DIRS=util session db csv ws server webui test
.PHONY: test

wreis:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	rm -rf dist

test:
	for dir in $(DIRS); do make -C $$dir test;done

package:
	for dir in $(DIRS); do make -C $$dir package;done

all: clean wreis package test
	echo "Completed"
