clean-data:
	rm -fr data/*

slave: clean-data
	go run main.go slave

.PHONY: slave clean-data