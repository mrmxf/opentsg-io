# set up initial variables of the github pat, the commit message and the target repo
git_pat=$1
message=$2
pullrepo=$3

targetRepo=https://mm-sandbox:$git_pat@github.com/mm-sandbox/$pullrepo.git



targetBranch="main"


# Import clog files
#source clogrc/privatepublic.sh
#. clogrc/privatepublic.sh

#update the names to reflect the new repo
# remove the old mod to stop changes
rm go.mod

# if statement here
#Walk_dir "$PWD" "gitlab.com/mmTristan/tpg-widgets" "github.com/mmTristan/tpg-widgets-private"
#Walk_dir "$PWD" "gitlab.com/mmTristan/tpg-core" "github.com/mmTristan/tpg-core-private"
#Walk_dir "$PWD" "gitlab.com/mmTristan/tpg-io" "github.com/mmTristan/tpg-io-private"


# make the new mod with the new name
go mod init "github.com/mmTristan/$pullrepo"
go mod tidy

# change these names to be whatever config
git config --global user.email "tristan@mrmxf.com"
git config --global user.name "mmTristan"

FULL=$(git log --format="%H" -n 1)
Branch=$(echo $FULL | head -c 38)


# make a temporary file to get the .git we want
# and replace the deleted one with it
mkdir tmp
cd tmp

if ! git clone --no-checkout $targetRepo ; then
	exit 1
fi



cd ..
# remove the original git
rm -R .git
# replace with the github target
cp -R "tmp/$pullrepo/.git" ".git" 
# remove the tmp folder
rm -R tmp


# add the branch to the hosted forked repo
git add . 
# remove uneeded pipeline files
git commit -m "$message"
git branch -M "$Branch"

if ! git push $targetRepo ; then
    echo "error pushing update to mm-sandbox/$pullrepo.git"
	exit 1
fi

#set authorisation here
export GH_TOKEN="$git_pat"
# pull request to the new library
if ! gh pr create --repo mmTristan/$pullrepo  --title "$message" --body "$message" --head "mm-sandbox:$Branch" --base $targetBranch ; then
    export GH_TOKEN=""
    echo "error generating pull request to mmTristan/$pullrepo.git"
	exit 1
fi
export GH_TOKEN=""

# remove the replaced git
rm -R .git