<!--- PROJECT LOGO --->
![project_logo](https://github.com/amcajal/8_bit_hubble/blob/master/project/doc/media/8_bit_hubble_golang_logo.png)

<!--- PROJECT SUMMARY/OVERVIEW --->
Reimplementation of **8_bit_hubble project** in Golang
Original project: https://github.com/amcajal/8_bit_hubble

Branches:
- master (entry point of the repository)
- scratch - all code started from zero using Golang
- replacement - certain C functionality is replaced (C-Golang bindings)

Replacement branch keeps most of the C code, but uses go to generate the PNG images

Core changes in the Golang version:
- All is implemented with Go.
- libpng is no longer used. All PNG operations are performed using the standard library (image package and so)
- Sprites are now handled in a completely different way (more details later)

In the original code, sprites where created as "CSV" files (using an Excel spreadsheet with format rules and so).
Then such CSV files were translated to C code using a Python script, and the C code must be added manually.

The Python script basically turned the CSV data representing the script into a sequence of operations in C.
Such operatons set bytes of a dynamic array to specific values (the array contains one byte per image pixel,
and its value is its color in hex format, i.e: 0xFF00FF)

Now, in Golang, the approach is different. Two ideas where tested:
- Host the png images containing the scripts in Github, and retrieve them on demand using the "net/http" package
- Turn the png images containing the scripts into base64 strings, and store them in the code itself as variables

The second approach is used.
- A png image containing a script (whatever it is) can be created with whatever software the user wants
- The png image is turned into a base64 string, using the golang executable provided to do so
- That string can be hardcoed into the code, and used when required

base64 string is decoded into a bytes slice, and such slice can be decoded into a "image.Image" variable.
And then, the sprite can be used in a very easy way.

This approach, besides exercising other parts of the golang language (which is the goal of this project),
allows to create sprites more easily, using whatever graphic editor the user needs. And working with "image.Image"
types allows to execute cool operations, like rescaling algorithms.

I.e: create a "planet" sprite with Gimp, save to base64, decode into image.Image, and reescale it 5 times to create
a Jupiter-alike sprite.
