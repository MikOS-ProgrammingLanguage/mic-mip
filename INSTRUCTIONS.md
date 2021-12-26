# Introduction to the mik compiler and the mip package manager

## Before use

Make sure to install via

    ./mic -install

At the location of the installation. Make sure to do this again everytime you move the compiler around in your filesystem.

## Use the mik compiler

If you have your first mik file ready, use

    ./mic -i your_file_name.mik

to compile it to standart C code. However if your target isn't C you can use other targets.

    ./mic -i your_file_name.mik -wasm

- You can try other languages as well. But not all will work. Right now three languages will get support and are availible with a command:

  1. **C**
  2. **Assembly (asm)**
  3. **Web assembly (wasm)**

But with no specified target C is set as standart.

## Use the mip package manager

Mip can be used for the four following tasks:

1. adding package structures to the **mik-src** folder
2. installing packages from github
3. removing installed packages
4. listing all installed packages

### 1 Adding package structures to mik-src

Lets say you have the following package structure:

    test_pkg/
        |____ testPkgFile1.milk
        |____ testPkgFile2.milk
        |____ milk.pkg

If you now want to 'yoink' the test package with

    #yoink-src <test_pkg>

you can add the package by using

    ./mip -add_pkg test_pkg

(use a relative or direct path)

You can now use

    #yoink-src <test_pkg>

to use all the files and functions in it.

### 2 Install packages from Github

If you want to use a package from Github, you can use

    sudo ./mip -install <github link>

to install it. You can now also use it with

    #yoink-src <pkg_name>

The pkg name needs to be provided by the repo or with the List command.

### 3 Remove a package

To remove a installed package just type

    sudo ./mip -remove <pkg_name>

This will remove the pkg from mik-src. Hence, you can't refference it with

    #yoink-src <pkg_name>

anymore.

### 4 List all packages

To list all packages type

    ./mip -list

The output will look something like

    You currently have 3 packages installed:
        |____ pkg_1
        |____ pkg_2
        |____ pkg_3
