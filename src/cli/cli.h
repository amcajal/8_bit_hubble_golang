/*******************************************************************************
*   Project: 8-Bit Hubble
*
*   File: cli.h
*
*   Description: Header file of CLI module. CLI module (Command Line
* 				Interface) implements the ENTRY POINT of the application
* 				in its CLI version. It parses the input arguments
* 				provided by the user, and performs the required
* 				operations (i.e: print version, help message, generate
* 				image, or print error message).
*
*   Notes: N/A
*
*   Contact: Alberto Martin Cajal, amartin.glimpse23<AT>gmail.com
*
*   URL: https://github.com/amcajal/8_bit_hubble
*
*   License: GNU GPL v3.0
*
*   Copyright (C) 2018 Alberto Martin Cajal
*
*   This file is part of 8-Bit Hubble.
*
*   8-Bit Hubble is free software: you can redistribute it and/or modify
*   it under the terms of the GNU General Public License as published by
*   the Free Software Foundation, either version 3 of the License, or
*   (at your option) any later version.
*
*   8-Bit Hubble is distributed in the hope that it will be useful,
*   but WITHOUT ANY WARRANTY; without even the implied warranty of
*   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*   GNU General Public License for more details.
*
*   You should have received a copy of the GNU General Public License
*   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*******************************************************************************/

#ifndef _8BH_CLI_H
#define _8BH_CLI_H

int cli_parser (int argc, char **argv);

#endif
