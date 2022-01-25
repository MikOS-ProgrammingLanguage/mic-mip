#!/bin/bash

function build {
    echo -e $(printf "\e[1;32mBuilding $1 ✨\e[0m")

    (
        cd $1 2> /dev/null
        go build 2> /dev/null
        mv $1 ../$2 2> /dev/null
    )
    if [ $? -ne 0 ];
    then
        echo -e $(printf "\e[31mFailed building $1 ❌\e[0m")
        # write error if requested
        #read var
        #if [ "$var" = "t" ]
        #then
        #    echo -e $(printf "\e[31mError:\e[0m")
        #    echo "$?"
        #fi
        exit 1
    else
        echo ✅
    fi
    
}


function install_ {
    echo -e $(printf "\e[1;32mInstalling $1 ✨\e[0m")

    (
        cp $1 /bin/$1 2> /dev/null
    )
    if [ $? -ne 0 ];
    then
        echo -e $(printf "\e[31mFailed install $1 ❌\e[0m")
        # write error if requested
        #read var
        #if [ "$var" = "t" ]
        #then
        #    echo -e $(printf "\e[31mError:\e[0m")
        #    echo "$?"
        #fi
        exit 1
    else
        echo ✅
    fi
    
}



if [ "$@" = "mic" ];
then
    build mic_ mic
fi

if [ "$@" = "mip" ];
then
    build mip_ mip
fi

if [ "$@" = "." ];
then
    build mic_ mic
    build mip_ mip
fi

if [ "$@" = "install" ];
then
    install_ mic
    install_ mip
fi
