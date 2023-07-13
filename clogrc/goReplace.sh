# this is designed to run strip a github module of it's private name
# so that it can be pushed as a fresh module


#Target="github.com/mmTristan/tpg-widgets-private"
#Update="github.com/mmTristan/tpg-widgets"


Walk_dir () {  
    #ignore empty and hidden folders
    shopt -s nullglob dotglob

    path="$1"
    Target="$2"
    Update="$3"
    

    for pathname in "$path"/*; do
        if [ -d "$pathname" ]; then # if it is a folder
           
            Walk_dir $pathname $Target $Update
        elif file_fence $pathname; then 
            #else update the targets with in the file
            echo "updating $Target to $Update in $pathname"
            sed -i "s~$Target~$Update~" $pathname
        fi
    done
}

file_fence () {

    goname=$1
    found=`echo $1| grep -e "\.go$"`
    foundmod=`echo $1| grep -e "go\.mod$"`
    #found just gives the matching string


    if [ "$found" = "$goname" ] || [ "$goname" = "$foundmod" ]; then
        return 0
    fi
    #return that no go files were indeed found
    return 1

}

#sed -i "s~$Target~$Update~" go.mod

#walk_dir $PWD

#go mod tidy

Walk_dir $1 $2 $3 