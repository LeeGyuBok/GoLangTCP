@echo off
REM This batch file generates Go code from the Protobuf definition

REM Display current directory
echo Current directory: %CD%

REM Run protoc command
echo Generating Go code from Protobuf...
protoc --csharp_out=. --grpc_out=. --plugin=protoc-gen-grpc="%GRPC%" *.proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto

REM Check if the command was successful
if %ERRORLEVEL% neq 0 (
    echo Error: Failed to generate Go code.
    exit /b 1
)

echo Go code generation completed successfully.

REM Pause to keep the command window open (optional)