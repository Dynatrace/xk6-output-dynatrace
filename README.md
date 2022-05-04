xk6-dynatrace-output
xk6.

Build
To build a k6 binary with this extension, first ensure you have the prerequisites:

Go toolchain
Git
Build with xk6:
xk6 build --with github.com/henrikrexed/xk6-dynatrace-output
This will result in a k6 binary in the current directory.

Run with the just build `k6:
./k6 run -o output-dynatrace <script.js>
