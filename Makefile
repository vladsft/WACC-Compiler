# Sample Makefile for the WACC Compiler lab: edit this to build your own compiler
# Locations

ANTLR_DIR   := antlr_config
SOURCE_DIR  := src
OUTPUT_DIR  := bin
INSTR_DIR   := src/assembly/instructions
TEMPLATE_DIR := src/assembly/instructions/templates

# Tools

ANTLR   := antlrBuild
FIND    := find
RM      := rm -rf
MKDIR   := mkdir -p
GO      := go

# the make rules

all: rules gen compile

# Generating the ANTLR parser and lexer files
# Compiling all the go packages in the source directory
rules: FORCE
	cd $(ANTLR_DIR) && ./$(ANTLR)

arm: src/assembly/instructions/arm.go

src/assembly/instructions/arm.go: $(TEMPLATE_DIR)/arm11_instructions.tmpl $(TEMPLATE_DIR)/asm_gen.go
	# cd $(TEMPLATE_DIR) && $(GO) run . arm11

gen: src/ast/acceptor.go src/visitor_generator/visitor.go
	# cd src && go generate ./...

compile: src/assembly/instructions/arm.go
	cd $(SOURCE_DIR) && $(GO) build -o ../compile

src/ast/acceptor.go: src/visitor_generator/visitor.go

src/visitor_generator/visitor.go:
	# cd src/visitor_generator && go run main.go

clean:
	$(RM) rules compile $(OUTPUT_DIR) $(SOURCE_DIR)/parser input input.s
	find src/ast/ -name "*visitor.go" -delete
	find src/ast/ -name "*acceptor.go" -delete
	find src/ -name "*_string.go" -delete

clean_assembly:
	find tests/ -name "*.s" -delete

FORCE:

.PHONY: all rules clean