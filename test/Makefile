DIRS = ws pic photo

clean:
	rm -f testreport.txt dbqqqmods.sql
	for dir in $(DIRS); do make -C $$dir clean;done
	@echo "*** CLEAN completed in test ***"

test:
	@echo "FUNCTIONAL TEST Report" > testreport.txt
	@echo "Report Start: \c" >> testreport.txt
	@date >> testreport.txt
	@echo >> testreport.txt
	@echo "Result  Path                  Test Name                                 Phases" >> testreport.txt
	@echo "------  --------------------  ----------------------------------------  ------" >> testreport.txt
	for dir in $(DIRS); do make -C $$dir test;done
	@./totals.sh
	@echo >> testreport.txt
	@echo "Report Complete: \c" >> testreport.txt
	@date >> testreport.txt
	@echo "*** TEST completed in test ***"
	@cat testreport.txt
	# Put usable data in the database...
	@mysql --no-defaults wreis < ws/xh.sql

package:
	for dir in $(DIRS); do make -C $$dir package;done
	@echo "*** PACKAGE completed in test ***"

secure:
	for dir in $(DIRS); do make -C $${dir} secure;done
	@rm -f config.json confdev.json confprod.json

testdb:
	cd ws;mysql --no-defaults wreis <xh.sql
