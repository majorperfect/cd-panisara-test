cert:
	cd cert; sh ./gen.sh; cd ..
.PHONY: gen clean server client test cert 