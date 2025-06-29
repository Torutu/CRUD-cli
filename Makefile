APP_NAME = library-cli
SRC = main.go

.PHONY: all run clean windows linux mac

all: windows

windows:
	go build -o $(APP_NAME).exe $(SRC)

linux:
	go build -o $(APP_NAME) $(SRC)

mac:
	go build -o $(APP_NAME) $(SRC)

run:
	go run $(SRC)

clean:
ifeq ($(OS),Windows_NT)
	del /F /Q $(APP_NAME).exe $(APP_NAME)
else
	rm -f $(APP_NAME) $(APP_NAME).exe
endif
