# go build task

OUT = gorem
CLC = go

.PHONY : build
build:
	$(CLC) build -o $(OUT)

.PHONY : clean
clean :
	$(RM) $(OUT)

.PHONY : run
run :
	./$(OUT) -addr=":8080"

.PHONY : help
help :
	@echo "Usage: make <task>"
	@echo " task: build clean"
