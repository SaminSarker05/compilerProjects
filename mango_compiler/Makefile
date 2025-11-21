TARGET = mango
LLVM_CONFIG = llvm-config-18

all: parser.cpp lexer.cpp
	g++ -std=c++17 \
		`$(LLVM_CONFIG) --cxxflags` \
		-fexceptions \
		-o $(TARGET) \
		main.cpp parser.cpp lexer.cpp generator.cpp \
		`$(LLVM_CONFIG) --ldflags --libs core executionengine interpreter mcjit native --system-libs`

parser.cpp: parser.y
	bison -d -o parser.cpp parser.y

lexer.cpp: lexer.l parser.cpp
	flex -o lexer.cpp lexer.l

clean:
	rm -f $(TARGET) parser.cpp parser.hpp lexer.cpp

test: all
	@echo "int x = 5;" | ./$(TARGET)
	@echo "int a = 2; int b = 3; int c = a + b;" | ./$(TARGET)

.PHONY: all clean test
