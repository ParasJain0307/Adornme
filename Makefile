# Makefile to generate Go server code & Swagger JSON for all components in one folder

COMPONENTS := Admin Carts Orders Payments Products Shipping Users
OUT_DIR=AdornmeCode

.PHONY: all clean $(COMPONENTS)

# Default: generate all
all: $(COMPONENTS)

# Generate server + Swagger for each component
$(COMPONENTS):
	@echo "Generating code for $@..."
	swagger generate server -f $(CURDIR)/$@/*.yaml -A adornme -t $(OUT_DIR)/$(shell echo $@ | tr '[:upper:]' '[:lower:]') --exclude-main
	swagger generate spec -f $(CURDIR)/$@/*.yaml -o $(OUT_DIR)/$(shell echo $@ | tr '[:upper:]' '[:lower:]')/swagger.json --scan-models

# Clean generated code
clean:
	@echo "Cleaning generated folder..."
	rm -rf $(OUT_DIR)/*

