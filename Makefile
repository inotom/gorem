# go build task

OUT = gorem
CLC = go

.PHONY : build
build:
	$(CLC) build -o $(OUT)

.PHONY : clean
clean :
	$(RM) $(OUT)

.PHONY : help
help :
	@echo "Usage: make <task>"
	@echo " task: build clean"
